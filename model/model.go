package model

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type DBHandler interface {
	AddUser(name, email, password string, authority int) *User
	GetUsers(userAuth int) []*User
	GetUser() *User
	ChangeUserAuth(userAuth int) *User
	DeleteUser(userAuth int) bool
	Close()
}
