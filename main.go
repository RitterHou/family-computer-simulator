package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)

	debug := flag.Bool("d", false, "Show debug message")
	flag.Parse()
	if *debug {
		fmt.Println("Debug mode opened...")
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}
}

func main() {
	romName := flag.Arg(0)
	if romName == "" {
		log.Fatal("no game rom file found")
	}

	rom := readFile(romName)
	log.Println(rom)
}
