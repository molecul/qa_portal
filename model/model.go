package model

import (
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/molecul/qa_portal/util/database"
)

type Configuration struct {
	LocalChallengesPath string
}

func syncLocalChallenge(localChallengePath string) error {
	return nil
}

func syncLocalChallenges(localChallengesPath string) error {
	files, err := ioutil.ReadDir(localChallengesPath)
	if err != nil {
		return fmt.Errorf("Error when traveling dir %v: %v", localChallengesPath, err)
	}

	for _, file := range files {
		logrus.Print(file)
	}

	return nil
}

func Init(cfg *Configuration) error {
	err := database.Get().Sync2(new(User), new(Test), new(Challenge))
	if err != nil {
		return err
	}

	return syncLocalChallenges(cfg.LocalChallengesPath)
}
