package main

import (
	"os/exec"
	"testing"
	"time"
)

// TODO test output returned
// Runs two jobs, one passes and one fails due to us stopping it
func TestJobRunner(t *testing.T) {
	// we'll run two commands, one will succeed and one will be stopped to simulate error
	job1, job2 := CreateJob(0, runSleepCmd+" 1"), CreateJob(1, runSleepCmd+" 1")
	jobsToRun, jobsCompleted, jobsErrored := make(chan *Job, 1), make(chan *Job, 1), make(chan *Job, 1)

	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)
	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)

	jobsToRun <- job2
	jobsToRun <- job1
	close(jobsToRun)

	<-time.After(100 * time.Millisecond) // wait for command to start before we stop it
	job2.Stop()                          // this one should error

	jobsDone := 0
	for jobsDone < 2 {
		select {
		case j := <-jobsCompleted:
			err := j.Cmd.Wait() // wait should be called by the runner before returning the job
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			}
		case j := <-jobsErrored:
			err := j.Cmd.Wait() // wait should be called by the runner before returning the job
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			} else if _, ok := j.Err.(*exec.ExitError); !ok {
				t.Fatalf("expected ExitError but got %v", err)
			}
		}

		jobsDone++
	}

}

var workingCmds = []string{"echo foo", "echo bar"}

// TODO test output
func TestCmdsImplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 0, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

// TODO test output
func TestCmdsExplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 1, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

// TODO test output
func TestCmdsParallelismHigherThanJobCount(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 100, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

func TestCmdsKeepOrder(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 0, true); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

var oneCmdFails = []string{"echo foo", "non-existent-command"}

func TestFailingCmds(t *testing.T) {
	if exitStatus := Cmds(oneCmdFails, 2, false); exitStatus == 0 {
		t.Fatalf("zero exit")
	}
}

// benchmarks

func BenchmarkCmdsEcho1(b *testing.B) {
	benchmarkCmdsEcho(1, b)
}

func BenchmarkCmdsEcho10(b *testing.B) {
	benchmarkCmdsEcho(10, b)
}

func BenchmarkCmdsEcho100(b *testing.B) {
	benchmarkCmdsEcho(100, b)
}
