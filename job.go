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

// a routine to run jobs from a channel until it closes
func jobRunner(jobsIn, jobsOut chan job) {
	for job := range jobsIn {
		out, err := job.cmd.CombinedOutput()
		job.out, job.err = string(out), err
		jobsOut <- job
	}
}
