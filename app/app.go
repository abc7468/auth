package app

import (
	"auth/api"
	"auth/auth"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 만약 access token, refresh token이 있다면
	// 이 때 이것이 유효하다면
	// http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)

	// 있지만 유효하지 않거나 없다면
	http.Redirect(w, r, "/login.html", http.StatusTemporaryRedirect)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
}

func MakeRouter() http.Handler {
	r := mux.NewRouter()
	api.AddApiRouter(r)
	auth.AddAuthRouter(r)
	n := negroni.Classic()
	n.UseHandler(r)
	r.HandleFunc("/", indexHandler).Methods("GET")
	return n
}
