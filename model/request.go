package model

type User struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Authority int    `json:"authority"`
}

type DataForToken struct {
	UserEmail      string `json:"user_email"`
	UserAuthorized string `json:"user_authorized"`
}
