package api

import (
	"auth/model"
	"auth/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var dbHandler model.DBHandler

type Test struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	test := &Test{}
	err := json.NewDecoder(r.Body).Decode(test)
	fmt.Println()
	fmt.Println(test)
	user, err := dbHandler.GetUser(test.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	fmt.Println(user)
	if ok := checkBcryptPassword(user.Password, test.Password); !ok {
		utils.RespondError(w, http.StatusUnauthorized, "incorrect password")
		return
	}

	data := map[string]string{"user_email": user.Email, "user_authorized": user.Authority}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/auth/token", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	utils.RespondJSON(w, http.StatusOK, res)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	inputEmail := r.FormValue("email")
	inputPwd := r.FormValue("password")
	inputName := r.FormValue("name")
	inputPhone := r.FormValue("phone")

	err := dbHandler.AddUser(inputName, inputEmail, string(generateBcryptPassword(inputPwd)), inputPhone)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
	}
	user := &model.User{Name: inputName, Email: inputEmail, Phone: inputPhone, Authority: "0"}
	utils.RespondJSON(w, http.StatusCreated, user)
}

func generateBcryptPassword(password string) []byte {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return bcryptPassword
}

func checkBcryptPassword(bcryptPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(bcryptPassword), []byte(inputPassword))
	return err == nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.NewRequest("DELETE", "http://localhost:8080/auth/token", nil)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := dbHandler.GetUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, users)
}
func AddApiRouter(r *mux.Router) {
	dbHandler = model.NewDBHandler()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", loginHandler).Methods("POST")
	apiRouter.HandleFunc("/signup", signupHandler).Methods("POST")
	apiRouter.HandleFunc("/logout", logoutHandler).Methods("POST")
	apiRouter.HandleFunc("/users", getUsersHandler).Methods("GET")
}
