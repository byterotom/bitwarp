// server.go
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
	*Node
}

func NewNodeServer(node *Node) *NodeServer {
	return &NodeServer{Node: node}
}

// function to run node server
func (nodeServer *NodeServer) Run(ready chan struct{}) {
	listner, err := net.Listen("tcp", "0.0.0.0:6969") 
	if err != nil {
		log.Fatalf("error initializing listner: %v", err)
	}

	grpcServer := grpc.NewServer()

	pbno.RegisterNodeServiceServer(grpcServer, nodeServer)

	serverPort := listner.Addr().(*net.TCPAddr).Port

	nodeServer.address = nodeServer.address + ":" + fmt.Sprint(serverPort)
	log.Printf("grpc node server running on %d...", serverPort)
	close(ready)
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
