package checker

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/molecul/qa_portal/model"
	"gopkg.in/yaml.v2"
)

type LocalChallengeSettings struct {
	Name       string
	Image      string
	Cmd        string
	Points     int64
	TargetPath string
	Inject     string
}

func loadYaml(path string, out interface{}) error {
	fileRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cannot open yaml file %s: %v", path, err)
	}
	if err = yaml.Unmarshal(fileRaw, out); err != nil {
		return fmt.Errorf("Cannot unmarshal yaml file %s: %v", path, err)
	}
	return nil
}

func (c *Checker) syncLocalChallenges(localChallengesPath string) error {
	files, err := ioutil.ReadDir(localChallengesPath)
	if err != nil {
		return fmt.Errorf("Error when parsing dir %v: %v", localChallengesPath, err)
	}

	for _, file := range files {
		if err := c.syncLocalChallenge(file.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Checker) loadLocalChallenge(directoryName string) (*model.Challenge, error) {
	dir := filepath.Join(c.Config.ChallengesPath, directoryName)

	settings := new(LocalChallengeSettings)
	if err := loadYaml(filepath.Join(dir, "settings.yml"), settings); err != nil {
		return nil, err
	}

	challenge := new(model.Challenge)
	challenge.InternalName = directoryName
	challenge.Image = settings.Image
	challenge.Name = settings.Name
	challenge.Points = settings.Points
	challenge.Cmd = settings.Cmd
	challenge.TargetPath = settings.TargetPath
	challenge.Inject = settings.Inject

	descPath := filepath.Join(dir, "README.md")
	description, err := ioutil.ReadFile(descPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot read %s file", descPath)
	}
	challenge.Description = string(description)

	return challenge, nil
}

func (c *Checker) syncLocalChallenge(directoryName string) error {
	logrus.Infof("Loading challenge %s", directoryName)
	chall, err := c.loadLocalChallenge(directoryName)
	if err != nil {
		return fmt.Errorf("Cannot load challenge %s:%v", directoryName, err)
	}

	existchall, err := model.GetChallengeByInternalName(directoryName)
	if err != nil {
		return err
	}

	if existchall == nil {
		return model.CreateChallenge(chall)
	} else {
		if !existchall.IsEqual(chall) {
			return existchall.UpdateWithInfoFrom(chall)
		} else {
			return nil
		}
	}
}
