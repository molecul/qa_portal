package checker

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/molecul/qa_portal/model"
)

type Configuration struct {
	Endpoint       string // docker socket path
	Timeout        int    // docker run timeout in milliseconds
	Threads        int    // docker runs in one time
	ChallengesPath string // path to etc/challenges folder
	ImagesPath     string // path to etc/challenges folder
}

type Checker struct {
	Docker *docker.Client
	Config *Configuration
}

var checker *Checker

func Get() *Checker {
	return checker
}

func Init(cfg *Configuration) (err error) {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "unix:///var/run/docker.sock"
	}

	c := &Checker{Config: cfg}
	if c.Docker, err = docker.NewClient(cfg.Endpoint); err != nil {
		return
	}
	checker = c
	return c.syncLocalChallenges(cfg.ChallengesPath)
}

func (c *Checker) CheckTest(challenge *model.Challenge, test *model.Test) {

}
