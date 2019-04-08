package main

import (
	"io/ioutil"
	"log"
)

func readFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
