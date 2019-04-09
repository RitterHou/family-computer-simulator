package main

import "testing"

func TestReadFile(t *testing.T) {
	t.Log(readFile("main.go"))
}

func TestNum(t *testing.T) {
	t.Log(0xff & 0xc >> 2)
}
