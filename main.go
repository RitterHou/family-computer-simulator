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
	rom := loadGameRom()
	fmt.Printf("%#v\n", rom)
}
