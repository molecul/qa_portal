package main

import (
	"context"
	"log"
	"time"

	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
)

func testChecker() {

}

func main() {
	c, err := checker.NewChecker(&checker.Settings{Endpoint: "tcp://127.0.0.1:6666"})
	if err != nil {
		panic(err)
	}

	task := c.NewTask(&model.Challenge{
		ID:            543,
		Image:         "python:2.7",
		ImageTestFile: "/tmp/task.py",
		ImageTestCmd:  "echo \"Inside $CHECKER_FILE:\"; cat $CHECKER_FILE",
		//ImageTestCmd: "for i in `seq 5`; do sleep 1; echo \"$i\"; done",
		InternalName: "test",
	}, &model.Test{
		ID:          1234,
		ChallengeID: 543,
		InputFile:   "1542 ololo 5f34 mda mda mda",
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err = task.Do(ctx); err != nil {
		log.Print(err)
	}
	cancel()

	log.Printf("ExitCode: %v", task.Result.ExitCode)
	log.Printf("Stdout:\n%v", task.Result.Stdout.String())
	log.Printf("Stderr:\n%v", task.Result.Stderr.String())
}
