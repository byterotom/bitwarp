package node

import (
	"log"

	pbno "github.com/Sp92535/internal/node/pb"
	pbtr "github.com/Sp92535/internal/tracker/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var trackerClient pbtr.TrackerServiceClient
var trackerConn *grpc.ClientConn

var nodeClient pbno.NodeServiceClient
var nodeConn *grpc.ClientConn

// function to declare node client
func NodeClientInit() {
	nodeConn, err := grpc.NewClient(":6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ERROR :%v", err)
	}

	nodeClient = pbno.NewNodeServiceClient(nodeConn)
}

// function to declare tracker client
func TrackerClientInit() {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	trackerConn, err := grpc.NewClient(":9999", opts)
	if err != nil {
		log.Fatalf("error connecting tracker :%v", err)
	}

	trackerClient = pbtr.NewTrackerServiceClient(trackerConn)
}

func StopNode() {
	trackerConn.Close()
	nodeConn.Close()
	// gracefull shutdown of node server yet to be added
}

func init() {
	TrackerClientInit()
	NodeClientInit()
}
