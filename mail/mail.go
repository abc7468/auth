package mail

import (
	"net/http"
	"net/smtp"

	"github.com/gorilla/mux"
)

func mailInit() {
	auth := smtp.PlainAuth("", "aqwer7468@gmail.com", "dd", "smtp.gmail.com")

	from := "aqwer7468@gmail.com"
	to := []string{"abc7468@naver.com"} // 복수 수신자 가능

	// 메시지 작성
	headerSubject := "Subject: 메일 인증\r\n"
	headerBlank := "\r\n"
	body := "<h1>[이메일 인증]</h1> <p>아래 링크를 클릭하시면 이메일 인증이 완료됩니다.</p> " +
		"<a href='http://localhost:8080/users/signup/confirm?key=" + "abcde" + "' target='_blenk'>이메일 인증 확인</a>"
	msg := []byte(headerSubject + headerBlank + body)

	// 메일 보내기
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		panic(err)
	}
}

func SendVerifyingMailHandler(w http.ResponseWriter, r *http.Request) {
	mailInit()
}

func AddMailRouter(r *mux.Router) {
	authRouter := r.PathPrefix("/mail").Subrouter()
	authRouter.HandleFunc("/verifying", SendVerifyingMailHandler).Methods("POST")
}
