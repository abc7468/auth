package app

import (
	"auth/api"
	"auth/auth"
	"auth/mail"
	"auth/model"
	"auth/utils"
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/home.html")
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/signup.html")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080/api/users")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
	}
	var users []*model.User
	json.NewDecoder(resp.Body).Decode(&users)
	t, _ := template.ParseFiles("./public/admin.html")
	t.Execute(w, users)
}

func MakeRouter() http.Handler {
	r := mux.NewRouter()
	api.AddApiRouter(r)
	auth.AddAuthRouter(r)
	mail.AddMailRouter(r)
	n := negroni.Classic()

	n.UseHandler(r)
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/main", homeHandler).Methods("GET")
	r.HandleFunc("/signup", signUpHandler).Methods("GET")
	r.HandleFunc("/admin", adminHandler).Methods("GET")
	return n
}
