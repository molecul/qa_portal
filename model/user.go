package model

import (
	"strings"
	"time"

	"github.com/molecul/qa_portal/util/database"
)

type User struct {
	ID            int64 `xorm:"pk autoincr"`
	Score         int64
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
	Email         string    `xorm:"unique"`
	EmailVerified bool
	Picture       string
	Name          string
}

func GetUserById(id int64) (*User, error) {
	u := new(User)
	has, err := database.Get().Id(id).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	email = strings.ToLower(email)
	u := &User{Email: email}
	has, err := database.Get().Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return u, nil
}

func CreateUser(u *User) error {
	u.Email = strings.ToLower(u.Email)

	_, err := database.Get().InsertOne(u)
	return err
}

func (u *User) Update() error {
	_, err := database.Get().Id(u.ID).AllCols().Update(u)
	return err
}
