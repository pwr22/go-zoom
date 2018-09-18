package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func runCmds(cmdStrs []string) int {
	jobsToRun := make(chan job, len(cmdStrs)) // enough to buffer all commands
	jobsCompleted := make(chan job, len(cmdStrs))

	cmdRunner := func() {
		for job := range jobsToRun {
			out, err := job.cmd.CombinedOutput()
			job.out, job.err = string(out), err
			jobsCompleted <- job
		}
	}

	numOfRunners := *parallelism
	if numOfRunners == 0 {
		numOfRunners = len(cmdStrs)
	}

	for n := 1; n <= numOfRunners; n++ { // start runners
		go cmdRunner()
	}

	jobs := make([]job, len(cmdStrs))
	for idx, cmdStr := range cmdStrs { // send the work
		job := createJob(cmdStr)
		jobs[idx] = job
		jobsToRun <- job
	}
	close(jobsToRun) // everything is queued up - this signals to runners when they should finish up

	doneCount, erroring, exitStatus := 0, false, 0
	for doneCount < len(cmdStrs) {
		job := <-jobsCompleted
		fmt.Print(job.out)

		if job.err != nil {
			if !erroring { // kill other processes after first error - then ignore the cascade - we only care about the first error
				fmt.Println("got an error so stopping all commands - further errors will be ignored")
				fmt.Printf("error: %v\n", job.err)

				if exitErr, ok := job.err.(*exec.ExitError); ok { // propagate exit status if we can
					exitStatus = exitErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() // might need to be made platform specific
				} else { // we don't have a status so just use a generic error
					exitStatus = 1
				}

				for _, job := range jobs {
					job.stop()
				}

				erroring = true
			}
		}

		doneCount++
	}

	return exitStatus
}
