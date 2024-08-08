package isolate

import (
	"fmt"
	"os/exec"
)

type isolator struct {
}

var isolateBinPath = "./pkg/bin/isolate/isolate"

func (is *isolator) init() error {
	return exec.Command(isolateBinPath, "--init").Run()
}

func (is *isolator) run(binName string) error {
	return exec.Command(isolateBinPath, "--run", binName).Run()
}

func (is *isolator) Exec(binPath string) (*exec.Cmd, error) {
	isolate := isolator{}

	if err := isolate.init(); err != nil {
		return nil, fmt.Errorf("init")
	}

	panic("not implemented")
	// dst, err := os.Create("")
	// io.Copy()
}
