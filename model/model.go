package model

import "github.com/molecul/qa_portal/util/database"

type Configuration struct{}

func Init(cfg *Configuration) error {
	return database.Get().Sync2(new(User), new(Test), new(Challenge))
}
