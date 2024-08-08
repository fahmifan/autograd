package podman

import (
	"bytes"
	"fmt"
	"os/exec"
)

type CPP struct {
	MountDir        string
	ProgramFileName string
}

type Result struct {
	stdout []byte
}

func (res *Result) Output() []byte {
	return res.stdout
}

func (cr *CPP) Run() (Result, error) {
	args := []string{
		"run", "--rm",
		"--memory=512m",
		"--cpus=1",
		"--network=none",
		"-v", fmt.Sprintf(`%s:/src`, cr.MountDir),
		"-w", "/src",
		"gcc:latest",
		"/bin/bash", "-c", fmt.Sprintf("g++ -o app %s && timeout 10s ./app", cr.ProgramFileName),
	}

	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("podman", args...)
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		return Result{}, fmt.Errorf("")
	}

	return Result{
		stdout: stdout.Bytes(),
	}, nil
}
