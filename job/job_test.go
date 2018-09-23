package job

import "testing"

func TestStopNil(t *testing.T) {
	var j Job
	j.Stop()
}
