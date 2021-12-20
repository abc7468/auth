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

func setValidToken(userId, userAuth string) (*model.AccessToken, error) {
	at := &model.AccessToken{}
	token, err := createToken(userId, userAuth)
	if err != nil {
		return nil, err
	}
	err = createAuth(userId, token)
	if err != nil {
		return nil, err
	}
	at.AccessToken = token.AccessToken
	return at, nil
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
	at, err := setValidToken(userData.UserId, userData.UserAuthorized)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, at)
}

func createToken(userId string, userAuth string) (*model.TokenDetails, error) {
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
	atClaims["user_id"] = userId
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
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("mcmvmkmsdnfsdmfdsjf"))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(userId string, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	errAccess := client.Set(td.AccessUuid, userId, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func createOnlyAccessToken(userId string, userAuth string, refreshTokenUuid string) (string, error) {
	expires := time.Now().Add(time.Minute * 15).Unix()
	accessUuid := uuid.NewV4().String()
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = userAuth
	atClaims["access_uuid"] = accessUuid
	atClaims["refresh_uuid"] = refreshTokenUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = expires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		return "", err
	}
	t := time.Unix(expires, 0)
	now := time.Now()
	errAccess := client.Set(accessUuid, userId, t.Sub(now)).Err()
	if errAccess != nil {
		return "", errAccess
	}
	return accessToken, nil
}
func CheckRemainTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "null" {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: false})
		return
	}
	token, err := ExtractTokenMetadata(r)

	if ok, _ := FetchAuthInRedis(token.AccessTokenUuid); ok {
		DeleteAuthForKey(token.AccessTokenUuid)
		DeleteAuthForKey(token.RefreshTokenUuid)
		at, err := setValidToken(token.UserId, token.UserAuthorized)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusOK, at)
		return
	}
	if ok, _ := FetchAuthInRedis(token.RefreshTokenUuid); ok {
		DeleteAuthForKey(token.RefreshTokenUuid)
		at, err := setValidToken(token.UserId, token.UserAuthorized)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusOK, at)
		return
	}
	if err != nil {
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: false})
		return
	}
	utils.RespondJSON(w, http.StatusOK, model.Success{Success: false})
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
		userId, ok := claims["user_id"].(string)
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
		tokenData.UserId = userId

		return tokenData, nil
	}
	return nil, err
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
	res, err := client.Get(uuid).Result()

	if err != nil {
		return false, err
	}
	if res == "" {
		return false, fmt.Errorf("no data in redis")
	}
	return true, nil
}

func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ExtractTokenMetadata(r)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	DeleteAuthForKey(token.AccessTokenUuid)
	DeleteAuthForKey(token.RefreshTokenUuid)
	utils.RespondJSON(w, http.StatusCreated, &model.Success{Success: true})
}

func DeleteAuthForKey(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

type Id struct {
	Id string `json:"id"`
}

func DeleteTokenForValueHandler(w http.ResponseWriter, r *http.Request) {
	userData := &model.UserId{}
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
	err = DeleteAuthForValue(userData.Id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, model.Success{Success: true})
}

func CheckTokenValidHandler(w http.ResponseWriter, r *http.Request) {
	returnToken := &model.AtAndSuccess{}
	returnToken.Success = false
	err := TokenValid(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, returnToken)
		return
	}
	token, err := ExtractTokenMetadata(r)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ok, _ := FetchAuthInRedis(token.AccessTokenUuid); !ok {
		if ok, _ := FetchAuthInRedis(token.RefreshTokenUuid); ok {
			at, err := createOnlyAccessToken(token.UserId, token.UserAuthorized, token.RefreshTokenUuid)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			returnToken.AccessToken = at
			utils.RespondJSON(w, http.StatusOK, returnToken)
			return
		} else {
			utils.RespondJSON(w, http.StatusOK, returnToken)
			return
		}
	}
	returnToken.Success = true
	utils.RespondJSON(w, http.StatusOK, returnToken)
}

func DeleteAuthForValue(userId string) error {
	keys, _, err := client.Scan(0, "", 0).Result()
	fmt.Println(err)
	for _, key := range keys {
		val, _ := client.Get(key).Result()
		if val == userId {
			DeleteAuthForKey(key)
		}
	}
	return nil
}

func getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ExtractTokenMetadata(r)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, data)
}

func createVerifyingCode(w http.ResponseWriter, r *http.Request) {
	//email을 key로
	email := &model.UserEmail{}
	utils.SetData(r, email)
	exp := time.Now().Add(time.Minute * 3).Unix()
	at := time.Unix(exp, 0) //converting Unix to UTC
	now := time.Now()
	code := uuid.NewV4().String()[:7]
	err := client.Set(email.Email, code, at.Sub(now)).Err()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
	}
	data := &model.VerificationData{}
	data.Email = email.Email
	data.Code = code
	utils.RespondJSON(w, http.StatusCreated, data)
}

func verifyingCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, _ := vars["code"]
	email := r.FormValue("email")
	val, err := client.Get(email).Result()
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, model.Success{Success: false})
		return
	}
	if code == val {
		client.Del(email)
		utils.RespondJSON(w, http.StatusOK, model.Success{Success: true})
		return
	}
	utils.RespondJSON(w, http.StatusUnauthorized, model.Success{Success: false})
}

func AddAuthRouter(r *mux.Router) {
	if client == nil {
		RedisConnect()
	}
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/token", CreateTokenHandler).Methods("POST")
	authRouter.HandleFunc("/token", CheckRemainTokenHandler).Methods("GET")
	authRouter.HandleFunc("/token/valid", CheckTokenValidHandler).Methods("GET")
	authRouter.HandleFunc("/token", DeleteTokenHandler).Methods("DELETE")
	authRouter.HandleFunc("/tokens", DeleteTokenForValueHandler).Methods("POST")
	authRouter.HandleFunc("/user", getUserDataHandler).Methods("GET")
	authRouter.HandleFunc("/code", createVerifyingCode).Methods("POST")
	//localhost:8080/auth/code/577bfbc?email=1@naver.com
	authRouter.Path("/code/{code}").Queries("email", "{email}").HandlerFunc(verifyingCode).Methods("GET")

}
