package checker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/molecul/qa_portal/model"
)

type Configuration struct {
	Endpoint       string        // docker socket path
	Timeout        time.Duration // docker run timeout in milliseconds
	Threads        int           // docker runs in one time
	ChallengesPath string        // path to etc/challenges folder
	ImagesPath     string        // path to etc/challenges folder
}

type Checker struct {
	Docker           *docker.Client
	Config           *Configuration
	Queue            chan *Task
	Close            chan struct{}
	CollectorUpdater chan bool
	testedNow        map[int64]bool
	testedMutex      sync.Mutex
}

var checker *Checker

func Get() *Checker {
	return checker
}

func (c *Checker) dispatch(task *Task) {
	c.testedMutex.Lock()
	if _, inProcess := c.testedNow[task.Test.Id]; !inProcess {
		c.testedNow[task.Test.Id] = true
		c.testedMutex.Unlock()
		defer func(id int64) {
			c.testedMutex.Lock()
			defer c.testedMutex.Unlock()
			delete(c.testedNow, id)
		}(task.Test.Id)
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*c.Config.Timeout)
		err := task.Do(ctx)
		if err != nil {
			log.Printf("[Checker-Dispatcher]: Error when execute task: %v", err)
		}
	} else {
		c.testedMutex.Unlock()
	}
}

func dispatcher(c *Checker) {
	for {
		select {
		case task := <-c.Queue:
			c.dispatch(task)
		case <-c.Close:
			return
		}
	}
}

func (c *Checker) clearCollectorChain() (times int) {
	for {
		select {
		case <-c.CollectorUpdater:
			times += 1
		default:
			return
		}
	}
}

func collector(c *Checker) {
	go func() {
		c.CollectorUpdater <- true
		ticker := time.Tick(time.Second * 10)
		for _ = range ticker {
			c.CollectorUpdater <- true
		}
	}()
	for {
		select {
		case <-c.CollectorUpdater:
			// Clear updater chan queue
			log.Printf("[Checker-Collector]: COLLECTOR UPDATING [%d]", c.clearCollectorChain())
			tests, err := model.TestsUntested(c.Config.Threads)
			if err != nil {
				log.Printf("[Checker-Collector]: Error loading tests: %v", err)
			}
			for _, test := range tests {
				c.CheckTest(test)
			}
		case <-c.Close:
			return
		}
	}
}

func Init(cfg *Configuration) (err error) {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "unix:///var/run/docker.sock"
	}

	c := &Checker{
		Config:           cfg,
		Queue:            make(chan *Task, cfg.Threads),
		CollectorUpdater: make(chan bool, 64),
		testedNow:        make(map[int64]bool),
		Close:            make(chan struct{}),
	}
	if c.Docker, err = docker.NewClient(cfg.Endpoint); err != nil {
		return
	}
	checker = c

	if err = c.PingDocker(); err != nil {
		return err
	}

	for i := 0; i < cfg.Threads; i++ {
		go dispatcher(c)
	}
	go collector(c)

	return c.syncLocalChallenges(cfg.ChallengesPath)
}

func (c *Checker) PingDocker() error {
	_, err := c.Docker.ListContainers(docker.ListContainersOptions{Limit: 1})
	return err
}

func (c *Checker) CheckTest(test *model.Test) {
	challenge, err := model.GetChallengeById(test.ChallengeId)
	if err != nil {
		log.Printf("[Checker-Collector]: Error loading challenge: %v", err)
	}
	task := c.NewTask(challenge, test)
	c.Queue <- task
}

func (c *Checker) Stop() {
	close(c.Close)
}
