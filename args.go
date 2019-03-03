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
// TODO pull out this parsing into a separate file and simplify it
func getArgSets() [][]string {
	var cmdSets [][]string

	inArgList, inFileList, setStartIdx := false, false, 0
	for i, a := range flag.Args() {
		// time to start building a new arg set
		if a == ":::" {
			// empty ::: or :::: are invalid
			if (inArgList || inFileList) && setStartIdx == i {
				var arg, symType string
				if inArgList {
					arg, symType = ":::", "arguments"
				} else {
					arg, symType = "::::", "files"
				}

				fmt.Fprintln(os.Stderr, arg+" must be followed by "+symType)
				os.Exit(2)
			} else if inFileList { // store the file set we were building
				for _, file := range flag.Args()[setStartIdx:i] {
					cmdSets = append(cmdSets, readCmdsFromFile(file))
				}
				inFileList = false
			} else if inArgList { // store the previous arg set we were building
				cmdSets = append(cmdSets, flag.Args()[setStartIdx:i])
			}

			// and now we are building a new arg set
			inArgList = true
			setStartIdx = i + 1
		} else if a == "::::" { // time to start building a new arg set
			// empty ::: or :::: are invalid
			if (inArgList || inFileList) && setStartIdx == i {
				var arg, symType string
				if inArgList {
					arg, symType = ":::", "arguments"
				} else {
					arg, symType = "::::", "files"
				}

				fmt.Fprintln(os.Stderr, arg+" must be followed by "+symType)
				os.Exit(2)
			} else if inArgList { // store the arg set we were building
				cmdSets = append(cmdSets, flag.Args()[setStartIdx:i])
				inArgList = false
			} else if inFileList { // store the previous file set we were building
				for _, file := range flag.Args()[setStartIdx:i] {
					cmdSets = append(cmdSets, readCmdsFromFile(file))
				}
			}

			// and now we're building a new file set
			inFileList = true
			setStartIdx = i + 1
		}
	}

	// trailing ::: or :::: are invalid
	if (inArgList || inFileList) && setStartIdx == len(flag.Args()) {
		var arg, symType string
		if inArgList {
			arg, symType = ":::", "arguments"
		} else {
			arg, symType = "::::", "files"
		}

		fmt.Fprintln(os.Stderr, arg+" must be followed by "+symType)
		os.Exit(2)
	} else if inArgList { // otherwise we just need to store the final set
		cmdSets = append(cmdSets, flag.Args()[setStartIdx:len(flag.Args())])
	} else if inFileList {
		for _, file := range flag.Args()[setStartIdx:len(flag.Args())] {
			cmdSets = append(cmdSets, readCmdsFromFile(file))
		}
	}

	return cmdSets
}

// return the command prefix given on the commandline, if any
func getCmdPrefix() string {
	start, end := 0, len(flag.Args())

	for i, a := range flag.Args() {
		if a == ":::" || a == "::::" { // stop as soon as we see any command arguments stuff
			end = i
			break
		}
	}

	return strings.Join(flag.Args()[start:end], " ")
}

// returns strings permuting all the argument sets given
func permuteArgSets(sets [][]string) []string {
	// no arg sets permutes to no args
	if len(sets) == 0 {
		return []string{}
	}

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

const placeHolder string = "{}"

// read in commands to run
func getCmdStrings() []string {
	prefix := getCmdPrefix()       // any command given as arguments on the command line
	cmdLineArgSets := getArgSets() // looks for ::: arguments
	var cmds []string

	// get commands
	if len(cmdLineArgSets) == 0 { // get args from stdin
		cmds = readCmdsFromFile("-")

	} else {
		cmds = permuteArgSets(cmdLineArgSets)
	}

	// build commands
	placeholderPresent := strings.Contains(prefix, placeHolder)
	for i, c := range cmds {
		if placeholderPresent {
			cmds[i] = strings.Replace(prefix, placeHolder, c, -1)
		} else { // default mode is to add arguments to the end of the command, if any
			cmds[i] = strings.Join([]string{prefix, c}, " ")
		}
	}

	return cmds
}
