package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func runCmds(cmdStrs []string, numOfRunners int) int {
	if numOfRunners == 0 { // default to running all commands at once
		numOfRunners = len(cmdStrs)
	}

	// start runners
	jobsToRun := make(chan job, len(cmdStrs)) // enough to buffer all commands
	jobsCompleted := make(chan job, len(cmdStrs))
	for n := 1; n <= numOfRunners; n++ {
		go jobRunner(jobsToRun, jobsCompleted)
	}

	// send out the work
	jobs := make([]job, len(cmdStrs))
	for idx, cmdStr := range cmdStrs { // send the work
		job := createJob(cmdStr)
		jobs[idx] = job
		jobsToRun <- job
	}
	close(jobsToRun) // everything is queued up - this signals to runners when they should finish up

	// receiving loop - waiting for jobs to come back from the runners
	doneCount, erroring, exitStatus := 0, false, 0
	for doneCount < len(cmdStrs) {
		job := <-jobsCompleted
		fmt.Print(job.out) // print out the output we got in all cases - success or failure

		if job.err != nil {
			if !erroring { // kill other processes after first error - then ignore the cascade - we only care about the first error
				fmt.Println("got an error so stopping all commands - further errors will be ignored")
				fmt.Printf("error: %v\n", job.err)

				if exitErr, ok := job.err.(*exec.ExitError); ok { // propagate exit status if we can
					exitStatus = exitErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() // might need to be made platform specific
				} else { // we don't have a status so just use a generic error
					exitStatus = 1
				}

				// stop all the running jobs
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
