// +build !windows

package job

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	job := Create("foo")

	if job.Err != nil {
		t.Fatal("err is not set to nil")
	}

	if job.Out != "" {
		t.Fatal("out is not empty")
	}

	if job.Cmd.SysProcAttr.Setpgid != true {
		t.Fatal("processes are not started in a new group")
	}

	if job.Cmd.Process != nil {
		t.Fatal("job should not have been started yet")
	}
}

func TestStopUnstarted(t *testing.T) {
	Create("sleep 0.1").Stop()
}

func TestStopStarted(t *testing.T) {
	job := Create("sleep 1")

	start := time.Now()
	if err := job.Cmd.Start(); err != nil {
		t.Fatal("could not start job")
	}

	job.Stop()
	err := job.Cmd.Wait()
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
