package main

import (
	"os/exec"
	"syscall"
)

// create a job to run a command
func createJob(cmdStr string) job {
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // create a new process group
	return job{cmd: cmd}
}

// stop a running job - no op if not running yet or already dead
func (job job) stop() {
	if job.cmd.Process != nil { // we can only do this if the process exists
		syscall.Kill(-job.cmd.Process.Pid, syscall.SIGTERM) // take down the process group
	}
}
