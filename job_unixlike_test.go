// +build !windows

package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestCreateJob(t *testing.T) {
	job := createJob("foo")

	if job.err != nil {
		t.Fatal("err is not set to nil")
	}

	if job.out != "" {
		t.Fatal("out is not empty")
	}

	if job.cmd.SysProcAttr.Setpgid != true {
		t.Fatal("processes are not started in a new group")
	}

	if job.cmd.Process != nil {
		t.Fatal("job should not have been started yet")
	}
}

func TestStopUnstartedJob(t *testing.T) {
	createJob("sleep 0.1").stop()
}

func TestStopStartedJob(t *testing.T) {
	job := createJob("sleep 1")

	start := time.Now()
	if err := job.cmd.Start(); err != nil {
		t.Fatal("could not start job")
	}

	job.stop()
	err := job.cmd.Wait()
	duration := time.Since(start)

	if err == nil {
		t.Fatal("should not be able to wait for process to finish")
	} else if err.Error() != "signal: terminated" {
		t.Fatalf("expected termination signal error but got: %v", err)
	}

	if int64(duration/time.Millisecond) >= 1000 {
		t.Fatal("command did not stop")
	}
}

// TODO test output returned
func TestJobRunner(t *testing.T) {
	// we'll run two commands, one will succeed and one will be stopped to simulate error
	job1, job2 := createJob("sleep 0.1"), createJob("sleep 1")
	jobsToRun, jobsCompleted, jobsErrored := make(chan job, 1), make(chan job, 1), make(chan job, 1)

	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)
	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)

	jobsToRun <- job2
	jobsToRun <- job1
	close(jobsToRun)

	<-time.After(100 * time.Millisecond) // wait for command to start before we stop it
	job2.stop()                          // this one should error

	jobsDone := 0
	for jobsDone < 2 {
		select {
		case job := <-jobsCompleted:
			err := job.cmd.Wait()
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			}
		case job := <-jobsErrored:
			err := job.cmd.Wait()
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			} else if _, ok := job.err.(*exec.ExitError); !ok {
				t.Fatalf("expected ExitError but got %v", err)
			}
		}

		jobsDone++
	}

}
