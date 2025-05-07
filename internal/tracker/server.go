package tracker

import (
	"fmt"
	"log"
	"net"

	"github.com/Sp92535/pkg"
	pbtr "github.com/Sp92535/proto/tracker/pb"
	"google.golang.org/grpc"
)

const PORT = 9999

type TrackerServer struct {
	pbtr.UnimplementedTrackerServiceServer
	address string
}

func NewTrackerServer() *TrackerServer {
	return &TrackerServer{
		address: pkg.GetLocalIp(),
	}
}

// function to run tracker server
func RunTrackerServer() {

	listner, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(PORT))
	if err != nil {
		log.Fatalf("error initializing listner: %v", err)
	}

	grpcServer := grpc.NewServer()
	trackerServer := NewTrackerServer()

	
	pbtr.RegisterTrackerServiceServer(grpcServer,trackerServer)
	go trackerServer.Sync()

	log.Printf("grpc tracker server running on %d...", PORT)
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
