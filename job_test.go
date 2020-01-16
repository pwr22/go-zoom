package main

import (
	"testing"
	"time"
)

func TestCreateJob(t *testing.T) {
	job := createJob(42, jobSleepCmd)

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
	var j job
	j.stop()
}

func TestStopUnstarted(t *testing.T) {
	createJob(42, jobSleepCmd).stop()
}

func TestStopStarted(t *testing.T) {
	job := createJob(42, jobSleepCmd)

	start := time.Now()
	if err := job.Cmd.Start(); err != nil {
		t.Fatal("could not start job")
	}

	job.stop()
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
		createJob(42, jobSleepCmd)
	}
}

func BenchmarkStopUnstartedJob(b *testing.B) {
	j := createJob(42, jobSleepCmd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.stop()
	}
}

func BenchmarkStartStopJob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := createJob(42, jobSleepCmd)

		if err := j.Cmd.Start(); err != nil {
			b.Fatal(err)
		}

		j.stop()

		if err := j.Cmd.Wait(); err == nil {
			b.Fatalf("expected an error")
		}
	}
}
