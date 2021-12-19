package main

import (
	"auth/app"
	"net/http"
)

func main() {
	a := app.MakeRouter()

	http.ListenAndServe(":8080", a)

}
