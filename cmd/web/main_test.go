package main

import "testing"

func TestRun(t *testing.T) {
	err, _ := initiate()
	if err != nil {
		t.Error("Failed initiate()")
	}
}
