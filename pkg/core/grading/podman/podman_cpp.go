package podman

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/fahmifan/autograd/pkg/core/grading"
)

type CPP struct{}

func (cr *CPP) Run(arg grading.RunnerArg) (grading.RunResult, error) {
	args := []string{
		"run", "--rm",
		"-i",
		fmt.Sprintf("--memory=%s", arg.MemLimit),
		"--network=none",
		"-v", fmt.Sprintf(`./%s:/src`, arg.MountDir),
		"-w", "/src",
		"docker.io/library/gcc:latest",
		"sh", "-c", fmt.Sprintf(`g++ -o /src/app /src/%s && timeout %s /src/app`, arg.ProgramFileName, arg.RunTimeout),
	}

	buffErr := bytes.NewBuffer(nil)

	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("podman", args...)

	buff, err := io.ReadAll(arg.Input)
	if err != nil {
		return grading.RunResult{}, fmt.Errorf("read input: %w", err)
	}
	cmd.Stdin = bytes.NewBuffer(buff)
	cmd.Stdout = stdout
	cmd.Stderr = buffErr

	if err := cmd.Run(); err != nil {
		return grading.RunResult{}, fmt.Errorf("run cpp: %w", err)
	}

	return grading.RunResult{
		Output: stdout.Bytes(),
	}, nil
}
