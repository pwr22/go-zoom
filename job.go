package main

import (
	"os/exec"
)

// a job to run a command
type job struct {
	cmd *exec.Cmd
	out string
	err error
}
