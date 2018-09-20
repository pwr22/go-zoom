package main

import (
	"fmt"
	"syscall"
	"testing"
)

func TestCreateJob(t *testing.T) {
	job := createJob("foo")

	if job.err != nil {
		t.Fatal("err is not set to nil")
	}

	if job.out != "" {
		t.Fatal("out is not empty")
	}

	if job.cmd.SysProcAttr.CreationFlags != syscall.CREATE_NEW_PROCESS_GROUP {
		t.Fatal("processes are not started in a new group")
	}

	if job.cmd.SysProcAttr.CmdLine != fmt.Sprintf(`/C "%s"`, "foo") {
		t.Fatal("CmdLine is not set correctly")
	}
}

func TestStop(t *testing.T) {
	job := createJob("ping -n 10 localhost")
	job.stop()
}
