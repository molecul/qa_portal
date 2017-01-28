package model

import (
	"strings"
	"time"

	"github.com/molecul/qa_portal/util/database"
	"github.com/molecul/qa_portal/util/isdebug"
)

type User struct {
	Id            int64     `xorm:"pk notnull autoincr"`
	Score         int64     `xorm:"notnull"`
	Created       time.Time `xorm:"notnull created"`
	Updated       time.Time `xorm:"notnull updated"`
	Email         string    `xorm:"notnull unique" json:"-"`
	EmailVerified bool      `json:"-"`
	Name          string    `xorm:"notnull"`
	Picture       string
}

func GetDebugUser() *User {
	return &User{
		Id:            0xDEADBEAF,
		Score:         0x31415234,
		Created:       time.Now().Add(-time.Hour * 24),
		Updated:       time.Now().Add(-time.Hour * 4),
		Email:         "debug_user@domain.com",
		EmailVerified: true,
		Name:          "Debug User",
	}
}

func GetUserById(id int64) (*User, error) {
	if isdebug.Is && GetDebugUser().Id == id {
		return GetDebugUser(), nil
	}
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
	if isdebug.Is && GetDebugUser().Email == email {
		return GetDebugUser(), nil
	}
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
	_, err := database.Get().Id(u.Id).AllCols().Update(u)
	return err
}

func Users(page, pageSize int, order string) ([]*User, error) {
	users := make([]*User, 0, pageSize)

	if order == "asc" {
		return users, database.Get().Limit(pageSize, (page-1)*pageSize).Asc("score").Find(&users)
	} else {
		return users, database.Get().Limit(pageSize, (page-1)*pageSize).Desc("score").Find(&users)
	}

}
