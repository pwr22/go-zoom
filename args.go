package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"
)

var printVersion = flag.BoolP("version", "V", false, "print version information")
var parallelism = flag.IntP("jobs", "j", runtime.NumCPU(), "number of jobs to run at once or 0 for as many as possible")
var keepOrder = flag.BoolP("keep-order", "k", false, "print output in the order jobs were run instead of the order they finish")
var dryRun = flag.Bool("dry-run", false, "print the commands that would be run instead of running them") // no shorthand to match GNU Parallel

// parse flags and commandline args
func parseArgs() {
	flag.SetInterspersed(false) // don't confuse flags to the command with our own
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *dryRun {
		for _, c := range getCmdStrings() {
			fmt.Println(c)
		}
		os.Exit(0)
	}
}

// return any command arguments given on the command line, if any
func getCmdLineArgs() [][]string {
	var cmdSets [][]string

	inList, start := false, 0
	for i, a := range flag.Args() {
		if a == ":::" {
			if !inList { // this is the first
				inList = true
			} else {
				cmdSets = append(cmdSets, flag.Args()[start:i])
			}

			start = i + 1 // we always restart counting
		}
	}

	if inList && start == len(flag.Args()) { // ::: was the last element
		fmt.Fprintln(os.Stderr, "::: must be followed by arguments")
		os.Exit(2)
	} else if inList {
		cmdSets = append(cmdSets, flag.Args()[start:len(flag.Args())])
	}

	return cmdSets
}

// return the command prefix given on the commandline, if any
func getCmdPrefix() string {
	start, end := 0, len(flag.Args())

	for i, a := range flag.Args() {
		if a == ":::" { // stop as soon as we see any command arguments stuff
			end = i
			break
		}
	}

	return strings.Join(flag.Args()[start:end], " ")
}

// returns strings permuting all the argument sets given
func permuteCmdLineArgSets(sets [][]string) []string {
	totalCmds := 1
	for _, cs := range sets {
		totalCmds *= len(cs)
	}

	cmds := make([]string, totalCmds)
	indices := make([]int, len(sets)) // tracks what permutation we're building
	for jobIdx := 0; jobIdx < totalCmds; jobIdx++ {
		cmdParts := make([]string, len(sets))

		for cs, i := range indices {
			cmdParts[cs] = sets[cs][i]
		}
		cmds[jobIdx] = strings.Join(cmdParts, " ")

		for cs := range indices {
			cs = len(indices) - 1 - cs

			if indices[cs] < len(sets[cs])-1 { // if we can do the next component in the current set we're done
				indices[cs]++
				break
			} else { // we need to consider the next component set
				indices[cs] = 0
			}
		}
	}

	return cmds
}

// reads in commands from a file with "-" meaning std in
func readCmdsFromFile(name string) []string {
	var file io.Reader
	if name == "-" {
		file = os.Stdin
	} else {
		f, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		file = f
	}

	// one command per line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// slurp in all the lines
	var cmds []string
	for scanner.Scan() {
		cmds = append(cmds, scanner.Text())
	}

	return cmds
}

// read in commands to run
func getCmdStrings() []string {
	prefix := getCmdPrefix()           // any command given as arguments on the command line
	cmdLineArgSets := getCmdLineArgs() // looks for ::: arguments
	var cmds []string

	// get commands
	if len(cmdLineArgSets) == 0 { // get args from stdin
		cmds = readCmdsFromFile("-")

	} else {
		cmds = permuteCmdLineArgSets(cmdLineArgSets)
	}

	// prefix commands
	for i, c := range cmds {
		cmds[i] = strings.Join([]string{prefix, c}, " ")
	}

	return cmds
}
