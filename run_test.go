package main

import "testing"

var workingCmds = []string{"echo foo", "echo bar"}

func TestRunCmdsImplicitParallelism(t *testing.T) {
	if exitStatus := runCmds(workingCmds, 0); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

func TestRunCmdsExplicitParallelism(t *testing.T) {
	if exitStatus := runCmds(workingCmds, 1); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

func TestRunCmdsParallelismTooHigh(t *testing.T) {
	if exitStatus := runCmds(workingCmds, 100); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}
