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
		log.Fatalf("invalid warp path")
	}

	n := node.NewNode(warp)
	go n.RunNodeServer()

	// initialize queue
	node.QueueInit()
	defer node.StopQueue()

	// initialize tracker client
	node.TrackerClientInit()
	defer node.StopNode()

	// consume message
	node.ConsumeMessage(n)

	n.SendResourceRequest()
	select {} // temporarily blocking
}
