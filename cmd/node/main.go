package main

import (
	"log"
	"os"

	"github.com/byterotom/internal/node"
	"github.com/byterotom/pkg"
	"github.com/byterotom/pkg/warpgen"
)

func main() {
	// reject if invalid method
	if len(os.Args) < 3 || (os.Args[1] != "seed" && os.Args[1] != "get") {
		log.Fatalf("specify valid method")
	}

	// reject if no warp file
	if len(os.Args) < 3 || os.Args[2] == "" {
		log.Fatalf("please specify warp file path")
	}

	isSeeder := os.Args[1] == "seed"

	// read warp file to warp
	warp := warpgen.ReadWarpFile(os.Args[2])
	if warp == nil {
		log.Printf("invalid warp path")
	}

	// ready signal channel to start register loop only when port number is appended to address
	ready := make(chan struct{})

	// intialize new node and node server
	n := node.NewNode(warp, isSeeder)
	nodeServer := node.NewNodeServer(n)
	defer node.StopNode()

	// run the node server
	go nodeServer.Run(ready)

	// start register loop once ready
	<-ready
	n.UpdateStatus() // this can update status in case of partial download

	// seeder doesnt download but will gracefully shutdown on CTRL+C
	if isSeeder {
		go n.RegisterLoop()
		select {} // temporarily blocking
	}

	d := pkg.RTT(n.Download)
	log.Printf("Downloaded in %f seconds", d)

	// merge all files once all chunks are here
	m := pkg.RTT(warp.MergeChunks)
	log.Printf("Merged in %f seconds", m)

}
