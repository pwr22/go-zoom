// +build !windows

package main

import (
	"syscall"
	"testing"
	"time"
)

const sleepCmd = "sleep"

var sleepCmds = []string{"sleep 1", "sleep 1"}

// unix only as windows doesn't have Kill
// TODO use a sigbreak on windows (it may not propagate)
func TestKeyboardInterruptCmds(t *testing.T) {
	time.AfterFunc(100*time.Millisecond, func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})

	if exitStatus := Cmds(sleepCmds, 0, false); exitStatus == 0 {
		t.Fatalf("zero exit")
	}
}
