package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
)

// a sensible default is to use the number of CPUs available
var parallelism = flag.Int("parallelism", runtime.NumCPU(), "how many commands to run at a time")

// parse flags and commandline args
func parseArgs() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "expecting a single argument - the path of the file of commands to run")
		os.Exit(2)
	}
}

// read in commands to run
func getCmdStrings() []string {
	cmdsFilePath := flag.Arg(0)

	cmdsFile, err := os.Open(cmdsFilePath)
	if err != nil {
		panic(err)
	}
	defer cmdsFile.Close()

	// one command per line
	scanner := bufio.NewScanner(cmdsFile)
	scanner.Split(bufio.ScanLines)

	// slurp in all the lines
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
