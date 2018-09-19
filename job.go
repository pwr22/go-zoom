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
func jobRunner(jobsIn, jobsFinished, jobsErrored chan job) {
	for job := range jobsIn {
		out, err := job.cmd.CombinedOutput()
		job.out, job.err = string(out), err
		if err != nil {
			jobsErrored <- job
		} else {
			jobsFinished <- job
		}
	}
}
