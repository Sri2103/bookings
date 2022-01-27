package main

import "testing"

func TestRun(t *testing.T) {

	err := Run()
	if err != nil {
		t.Errorf("testing failed for run()")
	}

}
