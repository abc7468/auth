package model

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Authority string `json:"authority"`
}

type DataForToken struct {
	AccessTokenUuid  string `json:"access_token_uuid"`
	RefreshTokenUuid string `json:"refresh_token_uuid"`
	UserId           string `json:"user_id"`
	UserAuthorized   string `json:"user_authorized"`
}

type Success struct {
	Success bool `json:"success"`
}

type VerificationData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
