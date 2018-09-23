package run

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/pwr22/zoom/job"
)

// jobRunner is a routine to run jobs from a channel until it closes
func jobRunner(jobsIn, jobsFinished, jobsErrored chan *job.Job) {
	for job := range jobsIn {
		out, err := job.Cmd.CombinedOutput()
		job.Out, job.Err = string(out), err
		if err != nil {
			jobsErrored <- job
		} else {
			jobsFinished <- job
		}
	}
}

// Cmds executes the commands its given in parallel
func Cmds(cmdStrs []string, numOfRunners int) int {
	if numOfRunners == 0 { // default to running all commands at once
		numOfRunners = len(cmdStrs)
	} else if numOfRunners > len(cmdStrs) { // or if there are more runners and commands then drop the excess
		numOfRunners = len(cmdStrs)
	}

	// start runners
	jobsToRun := make(chan *job.Job, len(cmdStrs)) // enough to buffer all commands
	jobsCompleted := make(chan *job.Job, len(cmdStrs))
	jobsErrored := make(chan *job.Job, len(cmdStrs))
	for n := 1; n <= numOfRunners; n++ {
		go jobRunner(jobsToRun, jobsCompleted, jobsErrored)
	}

	// send out initial jobs
	jobs := make([]*job.Job, len(cmdStrs))
	for idx := 0; idx < numOfRunners; idx++ {
		job := job.Create(cmdStrs[idx])
		jobs[idx] = job
		jobsToRun <- job
	}
	nextJobIdx := numOfRunners

	// if this is all the jobs then let the runners know before we even start collecting jobs
	if numOfRunners == len(cmdStrs) {
		close(jobsToRun)
	}

	stopEarlySignals := make(chan os.Signal)
	signal.Notify(stopEarlySignals, syscall.SIGINT, syscall.SIGTERM)

	// receiving loop - waiting for jobs to come back from the runners
	doneCount, stoppingEarly, exitStatus := 0, false, 0
	for doneCount < len(cmdStrs) {
		select {
		case finishedJob := <-jobsCompleted:
			fmt.Print(finishedJob.Out) // print out the output we got in all cases - success or failure

			if !stoppingEarly && nextJobIdx < len(cmdStrs) { // start any remaining jobs if things are still going smoothly
				nextJob := job.Create(cmdStrs[nextJobIdx])
				jobs[nextJobIdx] = nextJob
				jobsToRun <- nextJob
				nextJobIdx++

				// if there are no commands left then let the runners know
				if nextJobIdx == len(cmdStrs) {
					close(jobsToRun)
				}
			}

			doneCount++

		case erroredJob := <-jobsErrored:
			fmt.Print(erroredJob.Out) // print out the output we got in all cases - success or failure

			if !stoppingEarly { // kill other processes after first error - then ignore the cascade - we only care about the first error
				fmt.Println("got an error so stopping all commands - further errors will be ignored")
				fmt.Printf("error: %v\n", erroredJob.Err)

				if exitErr, ok := erroredJob.Err.(*exec.ExitError); ok { // propagate exit status if we can
					exitStatus = exitErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() // might need to be made platform specific
				} else { // we don't have a status so just use a generic error
					exitStatus = 1
				}

				// stop all the running jobs
				for _, j := range jobs {
					j.Stop()
				}

				stoppingEarly = true                   // make sure we don't come back into this branch
				doneCount += len(cmdStrs) - nextJobIdx // skip the jobs we aren't going to start
				if numOfRunners != len(cmdStrs) {      // safely handle the case where number of runners == number of jobs
					close(jobsToRun) // there will be no more work
				}
			}

			doneCount++

		case <-stopEarlySignals:
			if !stoppingEarly { // kill other processes after first error - then ignore the cascade - we only care about the first error
				fmt.Println("got a signal so stopping all commands")

				exitStatus = 1

				// stop all the running jobs
				for _, j := range jobs {
					j.Stop()
				}

				stoppingEarly = true                   // make sure we don't come back into this branch
				doneCount += len(cmdStrs) - nextJobIdx // skip the jobs we aren't going to start
				if numOfRunners != len(cmdStrs) {      // safely handle the case where number of runners == number of jobs
					close(jobsToRun) // there will be no more work
				}
			}
		}
	}

	return exitStatus
}
