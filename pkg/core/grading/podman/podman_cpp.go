package podman

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type Result struct {
	stdout []byte
}

func (res *Result) Output() []byte {
	return res.stdout
}

type CPP struct {
	MountDir        string
	ProgramFileName string
	Input           io.Reader
}

func (cr *CPP) Run() (Result, error) {
	args := []string{
		"run", "--rm",
		"-i",
		"--memory=100m",
		"--network=none",
		"-v", fmt.Sprintf(`./%s:/src`, cr.MountDir),
		"-w", "/src",
		"docker.io/library/gcc:latest",
		"sh", "-c", fmt.Sprintf(`g++ -o /src/app /src/%s && timeout 10s /src/app`, cr.ProgramFileName),
	}

	buffErr := bytes.NewBuffer(nil)

	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("podman", args...)

	buff, err := io.ReadAll(cr.Input)
	if err != nil {
		return Result{}, fmt.Errorf("read input: %w", err)
	}
	cmd.Stdin = bytes.NewBuffer(buff)
	cmd.Stdout = stdout
	cmd.Stderr = buffErr

	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("run cpp: %w", err)
	}

	return Result{
		stdout: stdout.Bytes(),
	}, nil
}
