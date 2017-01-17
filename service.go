package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/util/database"
	"github.com/molecul/qa_portal/web"
)

const TimestampFormat = "06-01-02 15:04:05.000"

type PrettyFormatter struct{}

func (f *PrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s|%7s|%s\n", entry.Time.Format(TimestampFormat), entry.Level.String(), entry.Message)), nil
}

func init() {
	logrus.SetFormatter(&PrettyFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	logrus.Infof("Starting server")

	if err := database.Init(&database.Configuration{
		Driver: "sqlite3",
		DSN:    "db.sqlite3",
	}); err != nil {
		logrus.Fatalf("Error when connecting to db: %v", err)
	}

	if err := model.Init(&model.Configuration{
		LocalChallengesPath: "etc/challenges",
	}); err != nil {
		logrus.Fatalf("Error syncing db: %v", err)
	}

	if err := checker.Init(&checker.Configuration{
		Endpoint: "tcp://127.0.0.1:6666",
	}); err != nil {
		logrus.Fatalf("Error when creating checker: %v", err)
	}

	web.Run(&web.Configuration{
		Hostname: "0.0.0.0",
		UseHTTP:  true,
		HTTPPort: 8000,
	})

	//test_checker(checker.Get())
}

func test_checker(c *checker.Checker) {
	task := c.NewTask(&model.Challenge{
		ID:         543,
		Image:      "python:2.7",
		TargetPath: "/tmp/task.py",
		Cmd:        "echo \"Inside $CHECKER_FILE:\"; cat $CHECKER_FILE",
		//Cmd: "for i in `seq 5`; do sleep 1; echo \"$i\"; done",
		InternalName: "test",
	}, &model.Test{
		ID:          1234,
		ChallengeID: 543,
		InputFile:   "1542 ololo 5f34 mda mda mda",
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := task.Do(ctx); err != nil {
		logrus.Print(err)
	}
	cancel()

	logrus.Printf("ExitCode: %v", task.Result.ExitCode)
	logrus.Printf("Stdout:\n%v", task.Result.Stdout.String())
	logrus.Printf("Stderr:\n%v", task.Result.Stderr.String())
}
