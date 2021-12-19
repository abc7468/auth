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
	AddUser(name, email, password, phone string) error
	GetUsers() ([]*User, error)
	GetUser(email string) (*User, error)
	ChangeUserAuth(userAuth, userId int) (*User, error)
	DeleteUser(userAuth, userId int) (bool, error)
	Close()
}

func NewDBHandler() DBHandler {
	return newMySQLHandler()
}
