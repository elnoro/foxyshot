package main

import (
	"foxyshot/cmd"
	"log"
	"os"
)

func main() {
	err := cmd.RunCmd(os.Args)
	if err != nil {
		log.Printf("Cannot run command, got error: %v", err)
	}
}
