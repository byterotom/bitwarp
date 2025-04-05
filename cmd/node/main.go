package main

import (
	"fmt"
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
		log.Fatalf("invalid path")
	}

	n := node.NewNode(warp)

	// initialize queue
	node.QueueInit()
	defer node.StopQueue()

	// initialize tracker client
	node.TrackerClientInit()
	// go node.RunNodeServer()
	defer node.StopNode()

	// consume message
	node.ConsumeMessage()

	n.SendResourceRequest()

	fmt.Println(n)

	select {} // temporarily blocking
}
