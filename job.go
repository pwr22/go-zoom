package main

import (
	"os/exec"
)

// Job to run a command
type Job struct {
	Cmd *exec.Cmd
	Out string // combined stdout / stderr
	Err error  // any error that occurred while running
	Num int    // the number of the command in the overall batch
}
