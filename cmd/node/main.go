package main

import (
	"log"
	"os"

	"github.com/Sp92535/internal/node"
	"github.com/Sp92535/pkg/warpgen"
)

func main() {
	// initialize queue
	node.QueueInit()
	defer node.StopQueue()

	// run tracker client
	node.RunTrackerClient()

	// consume message
	node.ConsumeMessage()

	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Fatalf("please specify warp file path")
	}

	warp := warpgen.ReadWarpFile(os.Args[1])
	log.Println(os.Args[1])
	if warp == nil{
		log.Fatalf("invalid path")
	}
	
	node.SendResourceRequest(warp.FileHash)

	select {} // temporarily blocking
}
