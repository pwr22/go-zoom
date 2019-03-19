// +build !windows

package main

import (
	"syscall"
	"testing"
	"time"
)

const runSleepCmd = "sleep"

var runSleepCmds = []string{"sleep 1", "sleep 1"}

// unix only as windows doesn't have Kill
// TODO use a sigbreak on windows (it may not propagate)
func TestKeyboardInterruptCmds(t *testing.T) {
	time.AfterFunc(100*time.Millisecond, func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})

	if exitStatus := Cmds(runSleepCmds, 0, false); exitStatus == 0 {
		t.Fatalf("zero exit")
	}
}

func benchmarkCmdsEcho(n int, b *testing.B) {
	// build a list of commands
	cmds := make([]string, n)
	for i := range cmds {
		cmds[i] = "echo foo > /dev/null"
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Cmds(cmds, 0, false)
	}
}
