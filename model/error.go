package model

import "fmt"

type ErrUserAlreadyExist struct {
	u *User
}

func (e ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("User already exist [id:%v, email:%v]", e.u.ID, e.u.Email)
}

func IsErrUserAlreadyExist(err error) bool {
	_, ok := err.(ErrUserAlreadyExist)
	return ok
}
