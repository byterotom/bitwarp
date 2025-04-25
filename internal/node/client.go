package node

import (
	"log"

	pbno "github.com/Sp92535/proto/node/pb"
	pbtr "github.com/Sp92535/proto/tracker/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var trackerClient pbtr.TrackerServiceClient
var trackerConn *grpc.ClientConn

// function to declare node client
func NodeClient(addr string) (*grpc.ClientConn, pbno.NodeServiceClient) {
	nodeConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ERROR :%v", err)
	}

	nodeClient := pbno.NewNodeServiceClient(nodeConn)
	return nodeConn, nodeClient
}

// function to declare tracker client
func TrackerClientInit() {
	var err error
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	trackerConn, err = grpc.NewClient("localhost:9999", opts)
	if err != nil {
		log.Fatalf("error connecting tracker :%v", err)
	}

	trackerClient = pbtr.NewTrackerServiceClient(trackerConn)
}

// function to clost tracker connction
func StopNode() {
	trackerConn.Close()
}

// initializer function (runs on any import of current file) to initialize tracker client
func init() {
	TrackerClientInit()
}
