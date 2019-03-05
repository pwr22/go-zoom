package main

import (
	"testing"
	"time"
)

func TestCreateJob(t *testing.T) {
	job := CreateJob(42, jobSleepCmd)

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
	CreateJob(42, jobSleepCmd).Stop()
}

func TestStopStarted(t *testing.T) {
	job := CreateJob(42, jobSleepCmd)

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

// benchmarking

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Create(42, sleepCmd)
	}
}

func BenchmarkStopUnstartedJob(b *testing.B) {
	j := Create(42, sleepCmd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.Stop()
	}
}

func BenchmarkStartStopJob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := Create(42, sleepCmd)

		if err := j.Cmd.Start(); err != nil {
			b.Fatal(err)
		}

		j.Stop()

		if err := j.Cmd.Wait(); err == nil {
			b.Fatalf("expected an error")
		}
	}
}
