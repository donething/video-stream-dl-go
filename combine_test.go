package main

import "testing"

func TestCombine(t *testing.T) {
	err := Combine(`D:\Data\Videos\20200901_105055`, ".ts", ".ts")
	if err != nil {
		t.Logf(err.Error())
	}
}
