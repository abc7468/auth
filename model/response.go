package model

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
type UserId struct {
	Id string `json:"id"`
}

type UserIdAndAuth struct {
	Id        string `json:"id"`
	Authority string `json:"authority"`
}

type AtAndSuccess struct {
	AccessToken string `json:"access_token"`
	Success     bool   `json:"success"`
}
