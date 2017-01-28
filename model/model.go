package model

import "github.com/molecul/qa_portal/util/database"

type Configuration struct {
	LocalTestFiles string
}

var configuration *Configuration

func Init(cfg *Configuration) error {
	configuration = cfg
	return database.Get().Sync2(new(User), new(Test), new(Challenge))
}
