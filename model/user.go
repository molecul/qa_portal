package model

import (
	"strings"
	"time"

	"github.com/molecul/qa_portal/util/database"
	"github.com/zalando/gin-oauth2/google"
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
	usr, err := GetUserByEmail(u.Email)
	if err != nil {
		return err
	}
	if usr != nil {
		return ErrUserAlreadyExist{u: usr}
	}

	_, err = database.Get().InsertOne(u)
	return err
}

func (u *User) Update() error {
	_, err := database.Get().Id(u.ID).AllCols().Update(u)
	return err
}

func (u *User) FillFromGoogle(gu *google.User) *User {
	u.Email = strings.ToLower(gu.Email)
	u.EmailVerified = gu.EmailVerified
	u.Name = gu.Name
	u.Picture = gu.Picture
	return u
}
