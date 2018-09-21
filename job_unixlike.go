// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

var shell = os.Getenv("SHELL")

// create a job to run a command
func createJob(cmdStr string) job {
	cmd := exec.Command(shell, "-c", cmdStr)              // assume the shell takes a command like this
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // create a new process group
	return job{cmd: cmd}
}

// stop a running job - no op if not running yet or already dead
func (job job) stop() {
	if job.cmd != nil && job.cmd.Process != nil { // we can only do this if the command and process exists
		syscall.Kill(-job.cmd.Process.Pid, syscall.SIGTERM) // take down the process group
	}
}
