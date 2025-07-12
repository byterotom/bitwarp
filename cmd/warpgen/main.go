package main

import (
	"log"
	"os"

	"github.com/byterotom/pkg/warpgen"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Fatalf("please provide file path.")
	}
	filePath := os.Args[1]
	warpgen.CreateWarpFile(filePath)
}
