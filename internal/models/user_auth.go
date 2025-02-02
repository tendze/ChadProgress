package models

// UserAuth used for
type UserAuth struct {
	Login    string
	Password string
}

func (u UserAuth) GetLogin() string {
	return u.Login
}

func (u UserAuth) GetPassword() string {
	return u.Password
}
