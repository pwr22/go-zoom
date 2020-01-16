// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

var shell = os.Getenv("SHELL")

// CreateJob returns a job to run a command.
func createJob(num int, cmdStr string) *job {
	cmd := exec.Command(shell, "-c", cmdStr)              // assume the shell takes a command like this
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // create a new process group
	return &job{Num: num, Cmd: cmd}
}

// Stop a running job - no op if not running yet or already dead.
func (j *job) stop() {
	if j != nil && j.Cmd != nil && j.Cmd.Process != nil { // we can only do this if a process exists
		syscall.Kill(-j.Cmd.Process.Pid, syscall.SIGTERM) // take down the process group
	}
}
