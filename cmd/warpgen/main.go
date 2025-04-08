package main

import (
	"log"
	"os"

	"github.com/Sp92535/pkg/warpgen"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("please provide file path.")
	}
	filePath := os.Args[1]
	warpgen.CreateWarp(filePath)
}
