package main

import (
	"fmt"
	"syscall"
	"testing"
)

const jobSleepCmd = "timeout 1"

func testCreateSpecificOS(t *testing.T, j *Job) {
	if j.Cmd.SysProcAttr.CreationFlags != syscall.CREATE_NEW_PROCESS_GROUP {
		t.Fatal("processes are not started in a new group")
	}

	if j.Cmd.SysProcAttr.CmdLine != fmt.Sprintf(`/C "%s"`, jobSleepCmd) {
		t.Fatal("CmdLine is not set correctly")
	}
}

func testStopErr(t *testing.T, err error) {
	// the error seems to vary on windows so we cannot test it meaningfully
	if err == nil {
		t.Fatal("should not be able to wait for process to finish")
	}
}
