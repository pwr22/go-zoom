package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"
)

var printVersion = flag.BoolP("version", "V", false, "print version information")
var parallelism = flag.IntP("jobs", "j", runtime.NumCPU(), "number of jobs to run at once or 0 for as many as possible")
var keepOrder = flag.BoolP("keep-order", "k", false, "print output in the order jobs were run instead of the order they finish")

// parse flags and commandline args
func parseArgs() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}
}

// read in commands to run
func getCmdStrings() []string {
	// get any command given as args
	commandPrefix := ""
	if len(flag.Args()) != 0 {
		commandPrefix = strings.Join(flag.Args(), " ")
	}

	// one command per line
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	// slurp in all the lines
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.Join([]string{commandPrefix, scanner.Text()}, " ")) // prefix with any commands
	}

	return lines
}
