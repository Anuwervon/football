package user

import "errors"

var Users []User
var idCnt int = 1

type User struct {
	ID       int
	Nickname string
	Login    string
	Password string
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
