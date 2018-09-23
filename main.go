package main

import (
	"os"

	"github.com/pwr22/zoom/run"
)

func main() {
	parseArgs()
	os.Exit(run.Cmds(getCmdStrings(), *parallelism))
}
