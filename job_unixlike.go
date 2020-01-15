// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

var shell = os.Getenv("SHELL")

// Create returns a job to run a command
func Create(num int, cmdStr string) *Job {
	cmd := exec.Command(shell, "-c", cmdStr)              // assume the shell takes a command like this
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // create a new process group
	return &Job{Num: num, Cmd: cmd}
}

// Stop a running job - no op if not running yet or already dead
func (job *Job) Stop() {
	if job != nil && job.Cmd != nil && job.Cmd.Process != nil { // we can only do this if a process exists
		syscall.Kill(-job.Cmd.Process.Pid, syscall.SIGTERM) // take down the process group
	}
}
