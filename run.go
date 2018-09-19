package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func runCmds(cmdStrs []string, numOfRunners int) int {
	if numOfRunners == 0 { // default to running all commands at once
		numOfRunners = len(cmdStrs)
	} else if numOfRunners > len(cmdStrs) { // or if there are more runners and commands then drop the excess
		numOfRunners = len(cmdStrs)
	}

	// start runners
	jobsToRun := make(chan job, len(cmdStrs)) // enough to buffer all commands
	jobsCompleted := make(chan job, len(cmdStrs))
	for n := 1; n <= numOfRunners; n++ {
		go jobRunner(jobsToRun, jobsCompleted)
	}

	// send out initial jobs
	jobs := make([]job, len(cmdStrs))
	for idx := 0; idx < numOfRunners; idx++ {
		job := createJob(cmdStrs[idx])
		jobs[idx] = job
		jobsToRun <- job
	}
	nextJobIdx := numOfRunners

	// if this is all the jobs then let the runners know before we even start collecting jobs
	if nextJobIdx == len(cmdStrs) {
		close(jobsToRun)
	}

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
				doneCount += len(cmdStrs) - nextJobIdx // skip the jobs we aren't going to start
				close(jobsToRun)
			}
		} else if !erroring && nextJobIdx < len(cmdStrs) { // start any remaining jobs if things are still going smoothly
			job := createJob(cmdStrs[nextJobIdx])
			jobs[nextJobIdx] = job
			jobsToRun <- job
			nextJobIdx++

			// if there are no commands left then let the runners know
			if nextJobIdx == len(cmdStrs) {
				close(jobsToRun)
			}
		}

		doneCount++
	}

	return exitStatus
}
