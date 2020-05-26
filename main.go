package main

import (
	"log"
	"os"

	"github.com/paraizofelipe/gosub/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
