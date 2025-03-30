package node

import (
	"log"
	"net"

	pbno "github.com/Sp92535/internal/node/pb"
	"google.golang.org/grpc"
)

type NodeServer struct {
	pbno.UnimplementedNodeServiceServer
}

func NewNodeServer() *NodeServer {
	return &NodeServer{}
}

func RunNodeServer() {
	listner, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("error initializing listner: %v", err)
	}

	grpcServer := grpc.NewServer()
	nodeServer := NewNodeServer()

	pbno.RegisterNodeServiceServer(grpcServer, nodeServer)

	log.Printf("GRPC node server running on 6969...")
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
