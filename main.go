package main

import (
	"fmt"
	"os"

	"github.com/pwr22/zoom/run"
)

const version = "v0.1.2"

func main() {
	if exitEarly, err := parseArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	} else if exitEarly {
		os.Exit(0)
	}

	cmds, err := getCmdStrings()
	if err != nil {
		os.Exit(2)
	}

	os.Exit(run.Cmds(cmds, *parallelism, *keepOrder))
}
