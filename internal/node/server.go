package node

import (
	"fmt"
	"log"
	"net"

	pbno "github.com/Sp92535/proto/node/pb"
	"google.golang.org/grpc"
)

type NodeServer struct {
	pbno.UnimplementedNodeServiceServer
}

func NewNodeServer() *NodeServer {
	return &NodeServer{}
}

// function to run node server
func (n *Node) RunNodeServer(ready chan struct{}) {
	listner, err := net.Listen("tcp", ":0") // dynamic port for os to choose available one
	if err != nil {
		log.Fatalf("error initializing listner: %v", err)
	}

	grpcServer := grpc.NewServer()
	nodeServer := NewNodeServer()

	pbno.RegisterNodeServiceServer(grpcServer, nodeServer)

	serverPort := listner.Addr().(*net.TCPAddr).Port

	n.address = n.address + ":" + fmt.Sprint(serverPort)
	log.Printf("GRPC server running on %d...", serverPort)
	close(ready)
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
