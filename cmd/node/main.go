package main

import (
	"log"
	"os"

	"github.com/Sp92535/internal/node"
	"github.com/Sp92535/pkg/warpgen"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Fatalf("please specify warp file path")
	}

	warp := warpgen.ReadWarpFile(os.Args[1])
	if warp == nil {
		log.Printf("invalid warp path")
	}

	// ready signal channel
	ready := make(chan struct{})

	n := node.NewNode(warp)
	go n.RunNodeServer(ready)

	<-ready
	n.Register()

	n.SendResourceRequest()

	defer node.StopNode()
	select {} // temporarily blocking

}
