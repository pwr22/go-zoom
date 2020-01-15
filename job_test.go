package main

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	job := Create(42, jobSleepCmd)

	if job.Err != nil {
		t.Fatal("err is not set to nil")
	}

	if job.Out != "" {
		t.Fatal("out is not empty")
	}

	if job.Num != 42 {
		t.Fatal("Num is not set")
	}

	if job.Cmd == nil {
		t.Fatal("Cmd is not set")
	}

	if job.Cmd.Path != shell {
		t.Fatal("The Cmd is not using $SHELL")
	}

	testCreateSpecificOS(t, job)

	if job.Cmd.Process != nil {
		t.Fatal("job should not have been started yet")
	}
}

func TestStopNil(t *testing.T) {
	var j Job
	j.Stop()
}

func TestStopUnstarted(t *testing.T) {
	Create(42, jobSleepCmd).Stop()
}

func TestStopStarted(t *testing.T) {
	job := Create(42, jobSleepCmd)

	start := time.Now()
	if err := job.Cmd.Start(); err != nil {
		t.Fatal("could not start job")
	}

	job.Stop()
	err := job.Cmd.Wait()
	duration := time.Since(start)

	testStopErr(t, err)

	if int64(duration/time.Millisecond) >= 1000 {
		t.Fatal("command did not stop")
	}
}
