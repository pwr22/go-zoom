package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// create a job to run a command
func createJob(cmdStr string) job {
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       fmt.Sprintf(`/C "%s"`, cmdStr),   // got to do weird things for the quoting to work right
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // signals go to the whole group so we gotta make a new one
	}
	return job{cmd: cmd}
}

// stop a running job - no op if not running yet or already dead
func (job job) stop() {
	if job.cmd != nil && job.cmd.Process != nil { // we can only do this if the command and process exists
		sendCtrlBreak(job.cmd.Process.Pid) // this goes to the whole process group and can only be sent within the same console
	}
}

// found in the tests for go itself
func sendCtrlBreak(pid int) {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		panic(err)
	}

	proc, err := dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		panic(err)
	}

	res, _, err := proc.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if res == 0 { // this seems to happen if the process is already dead so we can ignore it
		fmt.Fprintln(os.Stderr, err)
	}
}
