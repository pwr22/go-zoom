package run

import "testing"

var workingCmds = []string{"echo foo", "echo bar"}

// TODO test output
func TestCmdsImplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 0, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

// TODO test output
func TestCmdsExplicitParallelism(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 1, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}

// TODO test output
func TestCmdsParallelismHigherThanJobCount(t *testing.T) {
	if exitStatus := Cmds(workingCmds, 100, false); exitStatus != 0 {
		t.Fatalf("non-zero exit %d", exitStatus)
	}
}
