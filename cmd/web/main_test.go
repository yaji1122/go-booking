package main

import "testing"

func TestRun(t *testing.T) {
	err, _ := run()
	if err != nil {
		t.Error("Failed run()")
	}
}
