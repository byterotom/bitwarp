package node

import (
	"context"
	"log"
	"time"

	pbtr "github.com/Sp92535/internal/tracker/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var trackerClient pbtr.TrackerServiceClient
var trackerConn *grpc.ClientConn

func RunNodeClient() {
	conn, err := grpc.NewClient(":6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ERROR :%v", err)
	}
	defer conn.Close()

	// client := pbno.NewNodeServiceClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// res, err := client.Ping(ctx, &pbno.Empty{})
	// if err != nil {
	// 	log.Fatalf("Could not Ping: %v", err)
	// }

	// fmt.Println("Server Response:", res.Time)
}

func RunTrackerClient() {
	trackerConn, err := grpc.NewClient(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error connecting tracker :%v", err)
	}

	trackerClient = pbtr.NewTrackerServiceClient(trackerConn)
}

func StopNode() {
	trackerConn.Close()
}

func SendResourceRequest(msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := trackerClient.SendResourceRequest(ctx, &pbtr.ResourceRequest{Msg: msg})
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
}
