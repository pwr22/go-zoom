package run

import "testing"

var workingCmds = []string{"echo foo", "echo bar"}

func TestCmdsImplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 0, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

func TestCmdsExplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 1, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

func TestCmdsParallelismTooHigh(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 100, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}
