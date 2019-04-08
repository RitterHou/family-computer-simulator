package main

import "testing"

func TestReadFile(t *testing.T) {
	t.Log(readFile("main.go"))
}
