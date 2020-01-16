package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// runJobs is a routine to run jobs from a channel until it closes
func runJobs(jobsIn, jobsFinished, jobsErrored chan *job) {
	for j := range jobsIn {
		out, err := j.Cmd.CombinedOutput()
		j.Out, j.Err = string(out), err
		if err != nil {
			jobsErrored <- j
		} else {
			jobsFinished <- j
		}
	}
}

type runState struct {
	cmdStrs                                                              []string
	jobsToRun, jobsCompleted, jobsErrored                                chan *job
	jobs                                                                 []*job
	numOfRunners, doneCount, exitStatus, nextJobToRunIdx, totalCmdsCount int
	lastJobPrinted                                                       int // used with keepOrder to track how many we've printed
	stoppingEarly, keepOrder                                             bool
	jobsToPrint                                                          map[int]*job // put jobs in here when finished ready to print
}

func printJobs(finishedJob *job, state *runState) {
	if state.keepOrder {
		state.jobsToPrint[finishedJob.Num] = finishedJob

		// if we've just finished the next job to print
		if state.lastJobPrinted == finishedJob.Num-1 {
			// print jobs till we hit the end or are waiting on a job
			for idx := finishedJob.Num; idx < state.totalCmdsCount && state.jobsToPrint[idx] != nil; idx++ {
				fmt.Print(state.jobsToPrint[idx].Out)
				delete(state.jobsToPrint, idx)
				state.jobsToPrint[idx] = nil
				state.lastJobPrinted = idx
			}
		}
	} else { // output immediately
		fmt.Print(finishedJob.Out)
	}
}

func (state *runState) handleFinishedJob(finishedJob *job) {
	printJobs(finishedJob, state)
	state.jobs[finishedJob.Num] = nil

	if !state.stoppingEarly && state.nextJobToRunIdx < state.totalCmdsCount { // start any remaining jobs if things are still going smoothly
		nextJob := createJob(state.nextJobToRunIdx, state.cmdStrs[state.nextJobToRunIdx])
		state.jobs[state.nextJobToRunIdx] = nextJob
		state.jobsToRun <- nextJob
		state.nextJobToRunIdx++

		// if there are no commands left then let the runners know
		if state.nextJobToRunIdx == state.totalCmdsCount {
			close(state.jobsToRun)
		}
	}

	state.doneCount++
}

func (state *runState) handleErroredJob(erroredJob *job) {
	printJobs(erroredJob, state)

	if !state.stoppingEarly { // kill other processes after first error - then ignore the cascade - we only care about the first error
		fmt.Println("got an error so stopping all commands - further errors will be ignored")
		fmt.Printf("error: %v\n", erroredJob.Err)

		if exitErr, ok := erroredJob.Err.(*exec.ExitError); ok { // propagate exit status if we can
			state.exitStatus = exitErr.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() // might need to be made platform specific
		} else { // we don't have a status so just use a generic error
			state.exitStatus = 1
		}

		// stop all the running jobs
		for _, j := range state.jobs {
			if j != nil {
				j.stop()
			}
		}

		state.stoppingEarly = true                                      // make sure we don't come back into this branch
		state.doneCount += state.totalCmdsCount - state.nextJobToRunIdx // skip the jobs we aren't going to start
		if state.numOfRunners != state.totalCmdsCount {                 // safely handle the case where number of runners == number of jobs
			close(state.jobsToRun) // there will be no more work
		}
	}

	state.doneCount++
}

func (state *runState) handleStopEarlySignal() {
	if !state.stoppingEarly { // kill other processes after first error - then ignore the cascade - we only care about the first error
		fmt.Println("got a signal so stopping all commands")

		state.exitStatus = 1

		// stop all the running jobs
		for _, j := range state.jobs {
			if j != nil {
				j.stop()
			}
		}

		state.stoppingEarly = true                                      // make sure we don't come back into this branch
		state.doneCount += state.totalCmdsCount - state.nextJobToRunIdx // skip the jobs we aren't going to start
		if state.numOfRunners != state.totalCmdsCount {                 // safely handle the case where number of runners == number of jobs
			close(state.jobsToRun) // there will be no more work
		}
	}
}

// runCmds executes the commands its given in parallel
func runCmds(cmdStrs []string, numOfRunners int, keepOrder bool) (exitStatus int) {
	if numOfRunners == 0 { // means run everything at once
		numOfRunners = len(cmdStrs)
	} else if numOfRunners > len(cmdStrs) { // or if there are more runners than commands then drop the excess
		numOfRunners = len(cmdStrs)
	}

	// start runners
	jobsToRun := make(chan *job, len(cmdStrs)) // enough to buffer all commands
	jobsCompleted := make(chan *job, len(cmdStrs))
	jobsErrored := make(chan *job, len(cmdStrs))
	for n := 1; n <= numOfRunners; n++ {
		go runJobs(jobsToRun, jobsCompleted, jobsErrored)
	}

	// send out initial jobs
	jobs := make([]*job, len(cmdStrs))
	for idx := 0; idx < numOfRunners; idx++ {
		job := createJob(idx, cmdStrs[idx])
		jobs[idx] = job
		jobsToRun <- job
	}

	// if this is all the jobs then let the runners know before we even start collecting jobs
	if numOfRunners == len(cmdStrs) {
		close(jobsToRun)
	}

	stopEarlySignals := make(chan os.Signal)
	signal.Notify(stopEarlySignals, syscall.SIGINT, syscall.SIGTERM)

	state := &runState{
		numOfRunners:    numOfRunners,
		cmdStrs:         cmdStrs,
		totalCmdsCount:  len(cmdStrs),
		jobs:            jobs,
		jobsToRun:       jobsToRun,
		jobsCompleted:   jobsCompleted,
		jobsErrored:     jobsErrored,
		nextJobToRunIdx: numOfRunners,
		lastJobPrinted:  -1,
		keepOrder:       keepOrder,
		jobsToPrint:     make(map[int]*job),
	}

	// receiving loop - waiting for jobs to come back from the runners
	for state.doneCount < state.totalCmdsCount {
		select {
		case j := <-jobsCompleted:
			state.handleFinishedJob(j)

		case j := <-jobsErrored:
			state.handleErroredJob(j)

		case <-stopEarlySignals:
			state.handleStopEarlySignal()
		}
	}

	signal.Stop(stopEarlySignals)

	return state.exitStatus
}
