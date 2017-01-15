package checker

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/mogaika/qa_portal/model"
)

type Settings struct {
	Endpoint string // docker socket path
	Timeout  int    // docker run timeout in milliseconds
	Threads  int    // docker runs in one time
}

type Checker struct {
	Docker *docker.Client
	Cfg    *Settings
}

var checker *Checker = nil

func GetChecker() *Checker {
	return checker
}

func NewChecker(cfg *Settings) (c *Checker, err error) {
	if checker != nil {
		panic("Checker already created")
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "unix:///var/run/docker.sock"
	}

	c = &Checker{Cfg: cfg}

	if c.Docker, err = docker.NewClient(cfg.Endpoint); err != nil {
		return
	}

	checker = c
	return
}

func (c *Checker) CheckTest(challenge *model.Challenge, test *model.Test) {

}
