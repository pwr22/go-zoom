// +build !windows

package job

import (
	"testing"
)

const sleep = "sleep 1"

func testSysProcAttr(t *testing.T, job *Job) {
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
