package main

import (
	"auth/app"
	"net/http"
)

func main() {
	r := app.MakeRouter()

	http.ListenAndServe(":8080", r)

}
