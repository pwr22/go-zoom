package job

import (
	"os/exec"
)

// Job to run a command
type Job struct {
	Cmd *exec.Cmd
	Out string
	Err error
}
