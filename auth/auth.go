package auth

import (
	"auth/model"
	"auth/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	token, err := createToken(userData.UserEmail, utils.StringToInt(userData.UserAuthorized))
	fmt.Println("AA : " + userData.UserEmail)
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

func createToken(userEmail string, userAuth int) (*model.TokenDetails, error) {
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

func AddAuthRouter(r *mux.Router) {
	if client == nil {
		RedisConnect()
	}
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/token", CreateTokenHandler).Methods("POST")
}
