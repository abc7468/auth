package auth

import (
	"auth/model"
	"auth/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

var client *redis.Client

func RedisConnect() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	userData := &model.DataForToken{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(body, userData)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := createToken(userData.UserEmail, userData.UserAuthorized)
	at := &model.AccessToken{}
	if err != nil {
		fmt.Println(err)
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = createAuth(userData.UserEmail, token)
	if err != nil {
		fmt.Println(err)
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	at.AccessToken = token.AccessToken
	utils.RespondJSON(w, http.StatusCreated, at)
}

func createToken(userEmail string, userAuth string) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = userAuth
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["refresh_uuid"] = td.RefreshUuid
	atClaims["user_email"] = userEmail
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_email"] = userEmail
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("mcmvmkmsdnfsdmfdsjf"))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(userEmail string, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	errAccess := client.Set(td.AccessUuid, userEmail, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, userEmail, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CheckRemainTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "null" {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: false})
		return
	}
	token, err := ExtractTokenMetadata(r)

	if ok, _ := FetchAuthInRedis(token.AccessTokenUuid); ok {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: true})
		return
	}
	if ok, _ := FetchAuthInRedis(token.RefreshTokenUuid); ok {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: true})
		return
	}
	if err != nil {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: false})
		return
	}
}
func ExtractTokenMetadata(r *http.Request) (*model.DataForToken, error) {
	tokenData := &model.DataForToken{}
	token, err := VerifyToken(r)
	if err != nil {
		fmt.Println(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		authorized, ok := claims["authorized"].(string)
		if !ok {
			return nil, err
		}
		userEmail, ok := claims["user_email"].(string)
		if !ok {
			return nil, err
		}
		accessTokenUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		refreshTokenUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, err
		}
		tokenData.AccessTokenUuid = accessTokenUuid
		tokenData.RefreshTokenUuid = refreshTokenUuid
		tokenData.UserAuthorized = authorized
		tokenData.UserEmail = userEmail

		return tokenData, nil
	}
	return nil, err
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("jdnfksdmfksd"), nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("unvalid token")
	}
	return nil
}

func FetchAuthInRedis(uuid string) (bool, error) {
	_, err := client.Get(uuid).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ExtractTokenMetadata(r)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	DeleteAuth(token.AccessTokenUuid)
	DeleteAuth(token.RefreshTokenUuid)
	utils.RespondJSON(w, http.StatusCreated, &model.Success{Success: true})

}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func AddAuthRouter(r *mux.Router) {
	if client == nil {
		RedisConnect()
	}
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/token", CreateTokenHandler).Methods("POST")
	authRouter.HandleFunc("/token", CheckRemainTokenHandler).Methods("GET")
	authRouter.HandleFunc("/token", DeleteTokenHandler).Methods("DELETE")
}
