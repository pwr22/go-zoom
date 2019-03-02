// +build !windows

package run

import (
	"os/exec"
	"testing"
	"time"

	"github.com/pwr22/zoom/job"
)

// TODO test output returned
// Runs two jobs, one passes and one fails due to us stopping it
func TestJobRunner(t *testing.T) {
	// we'll run two commands, one will succeed and one will be stopped to simulate error
	job1, job2 := job.Create(0, "sleep 0.1"), job.Create(1, "sleep 1")
	jobsToRun, jobsCompleted, jobsErrored := make(chan *job.Job, 1), make(chan *job.Job, 1), make(chan *job.Job, 1)

	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)
	go jobRunner(jobsToRun, jobsCompleted, jobsErrored)

	jobsToRun <- job2
	jobsToRun <- job1
	close(jobsToRun)

	<-time.After(100 * time.Millisecond) // wait for command to start before we stop it
	job2.Stop()                          // this one should error

	jobsDone := 0
	for jobsDone < 2 {
		select {
		case j := <-jobsCompleted:
			err := j.Cmd.Wait() // wait should be called by the runner before returning the job
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			}
		case j := <-jobsErrored:
			err := j.Cmd.Wait() // wait should be called by the runner before returning the job
			if err == nil {
				t.Fatal("expected error that wait has already been called")
			} else if err.Error() != "exec: Wait was already called" {
				t.Fatalf("expected error that wait was already called but got: %v", err)
			} else if _, ok := j.Err.(*exec.ExitError); !ok {
				t.Fatalf("expected ExitError but got %v", err)
			}
		}

		jobsDone++
	}

}
