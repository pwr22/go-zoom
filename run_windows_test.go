package main

import "testing"

const runSleepCmd = "timeout"

func benchmarkRunCmdsEcho(n int, b *testing.B) {
	// build a list of commands
	cmds := make([]string, n)
	for i := range cmds {
		cmds[i] = "echo foo > NUL"
	}

	for i := 0; i < b.N; i++ {
		runCmds(cmds, 0, false)
	}
}
