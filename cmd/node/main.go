package main

import (
	"log"
	"os"

	"github.com/Sp92535/internal/node"
	"github.com/Sp92535/pkg/warpgen"
)

func main() {

	if len(os.Args) < 3 || (os.Args[1] != "seed" && os.Args[1] != "get") {
		log.Fatalf("specify valid method")
	}

	if len(os.Args) < 3 || os.Args[2] == "" {
		log.Fatalf("please specify warp file path")
	}

	isSeeder := os.Args[1] == "seed"

	warp := warpgen.ReadWarpFile(os.Args[2])
	if warp == nil {
		log.Printf("invalid warp path")
	}

	// ready signal channel
	ready := make(chan struct{})

	n := node.NewNode(warp, isSeeder)
	go n.RunNodeServer(ready)

	<-ready
	go n.RegisterLoop()

	if !isSeeder {
		n.SendResourceRequest()
	}

	defer node.StopNode()
	select {} // temporarily blocking

}
