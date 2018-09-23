package main

import "testing"

func TestStopNilJob(t *testing.T) {
	job{}.stop()
}
