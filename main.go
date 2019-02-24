package main

import (
	"os"

	"github.com/pwr22/zoom/run"
)

const version = "v0.1.1"

func main() {
	parseArgs()
	os.Exit(run.Cmds(getCmdStrings(), *parallelism, *keepOrder))
}
