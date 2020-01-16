// +build !windows

package main

import (
	"testing"
)

const jobSleepCmd = "sleep 1"

func testCreateSpecificOS(t *testing.T, j *job) {
	if len(j.Cmd.Args) != 3 || j.Cmd.Args[0] != shell || j.Cmd.Args[1] != "-c" || j.Cmd.Args[2] != jobSleepCmd {
		t.Fatal("The command to run is not set")
	}

	if j.Cmd.SysProcAttr.Setpgid != true {
		t.Fatal("processes are not started in a new group")
	}
}

func testStopErr(t *testing.T, err error) {
	if err == nil {
		t.Fatal("should not be able to wait for process to finish")
	} else if err.Error() != "signal: terminated" {
		t.Fatalf("expected termination signal error but got: %v", err)
	}
}
