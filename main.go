package main

import (
	"os"
)

func main() {
	parseArgs()
	os.Exit(runCmds(getCmdStrings(), *parallelism))
}
