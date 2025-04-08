package tracker

import (
	"fmt"
	"log"
	"net"

	pbtr "github.com/Sp92535/internal/tracker/pb"
	"google.golang.org/grpc"
)

const PORT = 9999

type TrackerServer struct {
	pbtr.UnimplementedTrackerServiceServer
}

func NewTrackerServer() *TrackerServer {
	return &TrackerServer{}
}

// function to run tracker server
func RunTrackerServer() {

	listner, err := net.Listen("tcp", ":"+fmt.Sprint(PORT))
	if err != nil {
		log.Fatalf("error initializing listner: %v", err)
	}

	grpcServer := grpc.NewServer()
	trackerServer := NewTrackerServer()

	pbtr.RegisterTrackerServiceServer(grpcServer,trackerServer)

	log.Printf("grpc tracker server running on %d...", PORT)
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
