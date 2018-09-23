package job

import (
	"fmt"
	"syscall"
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

	if job.Cmd.SysProcAttr.CreationFlags != syscall.CREATE_NEW_PROCESS_GROUP {
		t.Fatal("processes are not started in a new group")
	}

	if job.Cmd.SysProcAttr.CmdLine != fmt.Sprintf(`/C "%s"`, "foo") {
		t.Fatal("CmdLine is not set correctly")
	}
}

func TestStopUnstarted(t *testing.T) {
	Create("timeout 1").Stop()
}

func TestStopStarted(t *testing.T) {
	job := Create("timeout 1")

	start := time.Now()
	if err := job.Cmd.Start(); err != nil {
		t.Fatal("could not start job")
	}

	job.Stop()
	err := job.Cmd.Wait()
	duration := time.Since(start)

	if err == nil {
		t.Fatal("should not be able to wait for process to finish")
	}

	// the error seems to vary on windows so we cannot test it meaningfully

	if int64(duration/time.Millisecond) >= 1000 {
		t.Fatal("command did not stop")
	}
}
