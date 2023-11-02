package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/autograd/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
