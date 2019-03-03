// +build !windows

package job

import (
	"testing"
)

const sleepCmd = "sleep 1"

func testCreateSpecificOS(t *testing.T, job *Job) {
	if len(job.Cmd.Args) != 3 || job.Cmd.Args[0] != shell || job.Cmd.Args[1] != "-c" || job.Cmd.Args[2] != sleepCmd {
		t.Fatal("The command to run is not set")
	}

	if job.Cmd.SysProcAttr.Setpgid != true {
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
