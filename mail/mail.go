package mail

import (
	"auth/model"
	"auth/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/gorilla/mux"
)

var smtpAuth smtp.Auth

func initSmtpAuth() {
	smtpAuth = smtp.PlainAuth("", "aqwer7468@gmail.com", "rnfma12!", "smtp.gmail.com")

}

func mailInit(email *model.UserEmail) error {
	data := &model.UserEmailAndCode{}
	to := []string{email.Email} // 복수 수신자 가능
	from := "aqwer7468@gmail.com"

	jsonData, _ := json.Marshal(email)
	// 메시지 작성
	fmt.Println(email)

	resp, err := http.Post("http://localhost:8080/auth/code", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	json.NewDecoder(resp.Body).Decode(data)
	defer resp.Body.Close()
	subject := "Subject: Test email from Go!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<h1>[이메일 인증]</h1> <p>아래 코드를 입력하시면 이메일 인증이 완료됩니다.</p> " +
		"<h2>코드 : " + data.Code + "</h2>"
	msg := []byte(subject + mime + body)
	// 메일 보내기
	err = smtp.SendMail("smtp.gmail.com:587", smtpAuth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func SendVerifyingMailHandler(w http.ResponseWriter, r *http.Request) {
	email := &model.UserEmail{}
	utils.SetData(r, email)
	fmt.Println(email)
	err := mailInit(email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
	}
	utils.RespondJSON(w, http.StatusOK, model.Success{Success: true})
}

func AddMailRouter(r *mux.Router) {
	if smtpAuth == nil {
		initSmtpAuth()
	}
	authRouter := r.PathPrefix("/mail").Subrouter()
	authRouter.HandleFunc("/verifying", SendVerifyingMailHandler).Methods("POST")
}
