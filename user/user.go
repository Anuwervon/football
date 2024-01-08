package user

import "errors"

var Users []User
var idCnt int = 1

type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(u *User) error {
	ok := IsRegistered(*u)
	if ok {
		return errors.New("already registered")
	}

	u.ID = idCnt
	idCnt++
	Users = append(Users, *u)

	return nil
}

func IsRegistered(u User) bool {
	for _, v := range Users {
		if u.Login == v.Login {
			return true
		}
	}
	return false
}

func Login(u User) bool {
	for _, v := range Users {
		if v.Login == u.Login && v.Password == u.Password {
			return true
		}
	}

	return false
}

func FillUserFromLoginAndPassword(u *User) {
	for _, v := range Users {
		if v.Login == u.Login && v.Password == u.Password {
			u.ID = v.ID
			u.Nickname = v.Nickname
			return
		}
	}
}
