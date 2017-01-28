package checker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/molecul/qa_portal/model"
)

type Task struct {
	Docker    *docker.Client
	Test      *model.Test
	Challenge *model.Challenge
	Container *docker.Container
	Result    struct {
		ExitCode int
		Stdout   bytes.Buffer
		Stderr   bytes.Buffer
	}
}

func (c *Checker) NewTask(challenge *model.Challenge, test *model.Test) *Task {
	return &Task{
		Docker:    c.Docker,
		Test:      test,
		Challenge: challenge,
	}
}

func (task *Task) dockerInjectFiles(fnameDataMap map[string][]byte) error {
	tarBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuf)
	for fname, fdata := range fnameDataMap {
		if err := tw.WriteHeader(&tar.Header{
			Name: fname,
			Mode: 0644,
			Size: int64(len(fdata)),
		}); err != nil {
			return fmt.Errorf("Error writing file '%s' header to tar archive: %v", fname, err)
		}

		if _, err := tw.Write(fdata); err != nil {
			return fmt.Errorf("Error writing file '%s' to tar archive: %v", fname, err)
		}
	}

	if err := task.Docker.UploadToContainer(task.Container.ID, docker.UploadToContainerOptions{
		InputStream: tarBuf,
		Path:        "/",
	}); err != nil {
		return fmt.Errorf("Error uploading tar archive to container: %v", err)
	}
	return nil
}

func (task *Task) dockerInjectFile(targetFileName string, data []byte) error {
	return task.dockerInjectFiles(map[string][]byte{targetFileName: data})
}

func (task *Task) injectChallengeFiles() error {
	tarBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuf)
	for _, inject := range strings.Split(task.Challenge.Inject, ",") {
		args := strings.Split(inject, ":")
		if len(args) != 2 {
			continue
		}
		inFile := filepath.Join(Get().Config.ChallengesPath, task.Challenge.InternalName, args[0])

		logrus.Printf("Injecting %s at %s", inFile, args[1])
		frdr, err := os.Open(inFile)
		if err != nil {
			return err
		}
		defer frdr.Close()
		finfo, err := frdr.Stat()
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(&tar.Header{
			Name: args[1],
			Mode: int64(finfo.Mode()),
			Size: finfo.Size(),
		}); err != nil {
			return fmt.Errorf("Error writing file '%s' header to tar archive: %v", inFile, err)
		}

		if _, err := io.Copy(tw, frdr); err != nil {
			return fmt.Errorf("Error writing file '%s' data to tar archive: %v", inFile, err)
		}
		frdr.Close()
	}
	if err := task.Docker.UploadToContainer(task.Container.ID, docker.UploadToContainerOptions{
		InputStream: tarBuf,
		Path:        "/",
	}); err != nil {
		return fmt.Errorf("Error uploading tar archive to container: %v", err)
	}
	return nil
}

func (task *Task) dockerLoadLogs() error {
	if err := task.Docker.Logs(docker.LogsOptions{
		Container:    task.Container.ID,
		Stdout:       true,
		Stderr:       true,
		OutputStream: &task.Result.Stdout,
		ErrorStream:  &task.Result.Stderr,
		RawTerminal:  true,
	}); err != nil {
		return fmt.Errorf("Error collecting logs from container: %v", err)
	}
	return nil
}

func (task *Task) runTask(ctx context.Context, duration *time.Duration) (err error) {
	createOptions := docker.CreateContainerOptions{
		Name: fmt.Sprintf("qachecker-%v-%x-%x", task.Test.Id, time.Now().UnixNano(), rand.Int31()),
		Config: &docker.Config{
			Tty:             true,
			NetworkDisabled: true,
			Image:           task.Challenge.Image,
			Env: []string{
				"CHECKER_CHALLENGE=" + task.Challenge.InternalName,
				"CHECKER_TEST=" + strconv.FormatInt(task.Test.Id, 10),
				"CHECKER_FILE=" + task.Challenge.TargetPath,
			},
		},
		HostConfig: &docker.HostConfig{
		// With this we cannot collect logs
		// AutoRemove: true,
		},
	}

	log.Print(task.Challenge.Cmd)
	if task.Challenge.Cmd != "" {
		createOptions.Config.Entrypoint = []string{"/bin/bash", "-c"}
		createOptions.Config.Cmd = []string{task.Challenge.Cmd}
	}

	if task.Container, err = task.Docker.CreateContainer(createOptions); err != nil {
		return fmt.Errorf("Error creating container: %v", err)
	}

	defer func() {
		logrus.Infof("autoremove %s", task.Container.ID)
		task.Docker.RemoveContainer(docker.RemoveContainerOptions{
			ID:    task.Container.ID,
			Force: true,
		})
	}()

	absPath, err := filepath.Abs(task.Test.GetInputFileName())
	if err != nil {
		return fmt.Errorf("Error generating abspath: %v", err)
	}

	testInput, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("Error reading input file %s: %v", task.Test.GetInputFileName(), err)
	}

	if err = task.dockerInjectFile(task.Challenge.TargetPath, testInput); err != nil {
		return err
	}

	if err = task.injectChallengeFiles(); err != nil {
		return fmt.Errorf("Error when trying to inject challenge files: %v", err)
	}

	startTime := time.Now()
	if err = task.Docker.StartContainerWithContext(task.Container.ID, &docker.HostConfig{}, ctx); err != nil {
		return fmt.Errorf("Error when starting docker container: %v", err)
	}

	task.Result.ExitCode, err = task.Docker.WaitContainerWithContext(task.Container.ID, ctx)
	if err != nil {
		task.dockerLoadLogs()
		return fmt.Errorf("Error wait container: %v", err)
	}
	if duration != nil {
		*duration = time.Now().Sub(startTime)
	}

	return task.dockerLoadLogs()
}

func (task *Task) Do(ctx context.Context) (err error) {
	test := task.Test

	err = task.runTask(ctx, &test.Duration)
	log.Printf("Task ended: %v", err)
	checkTime := time.Now()
	test.Checked = &checkTime
	test.IsSucess = err == nil && task.Result.ExitCode == 0

	var output []byte
	if err == nil {
		output = task.Result.Stdout.Bytes()
	} else {
		output = []byte(err.Error())
	}

	if err := test.Update(output); err != nil {
		return fmt.Errorf("Error when updating test: %v", err)
	}

	return nil
}
