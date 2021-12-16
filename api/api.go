package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// loginId := r.FormValue("id")
	// loginPwd := r.FormValue("password")
	// 인증서버로 token확인

	// if access token 인증 됐다면 main 화면으로

	// else if access token 인증 안됐고 refresh token 인증 됐다면 access token 생성 후 main 화면으로

	// else 아무것도 없다면 id pwd 확인 후 token 생성 후 main 화면으로
	data := map[string]string{"user_email": "abc7468@naver.com", "user_authorized": "1"}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/auth/token", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)
}

func AddApiRouter(r *mux.Router) {
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", loginHandler)
}
