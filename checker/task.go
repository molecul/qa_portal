package checker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
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

func (task *Task) runTask(ctx context.Context) (err error) {
	createOptions := docker.CreateContainerOptions{
		Name: fmt.Sprintf("qachecker-%v", task.Test.ID),
		Config: &docker.Config{
			Tty:             true,
			NetworkDisabled: true,
			Image:           task.Challenge.Image,
			Entrypoint:      []string{"/bin/bash", "-c"},
			Env: []string{
				"CHECKER_CHALLENGE=" + task.Challenge.InternalName,
				"CHECKER_TEST=" + strconv.FormatInt(task.Test.ID, 10),
				"CHECKER_FILE=" + task.Challenge.TargetPath,
			},
		},
		HostConfig: &docker.HostConfig{
		// TODO We really need this?
		// AutoRemove: true,
		},
	}
	if task.Challenge.Cmd != "" {
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

	// TODO Change input file to real file data
	if err = task.dockerInjectFile(task.Challenge.TargetPath, []byte(task.Test.InputFile)); err != nil {
		return err
	}

	if err = task.Docker.StartContainerWithContext(task.Container.ID, &docker.HostConfig{}, ctx); err != nil {
		return fmt.Errorf("Error when starting docker container: %v", err)
	}

	task.Result.ExitCode, err = task.Docker.WaitContainerWithContext(task.Container.ID, ctx)
	if err != nil {
		task.dockerLoadLogs()
		return fmt.Errorf("Error wait container: %v", err)
	}

	return task.dockerLoadLogs()
}

func (task *Task) Do(ctx context.Context) (err error) {
	err = task.runTask(ctx)

	return err
}
