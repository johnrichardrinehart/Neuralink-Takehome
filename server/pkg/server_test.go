package server

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"

	"google.golang.org/grpc"
)

var client pb.NLImageServiceClient
var testHostPort = "localhost:2223"

func TestMain(m *testing.M) {
	// start server
	lis, err := net.Listen("tcp", testHostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterNLImageServiceServer(s, &Server{
		Debug: true,
	})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// connect client
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, testHostPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	client = pb.NewNLImageServiceClient(conn)

	// run tests
	exitCode := m.Run()

	// clean up
	if err := conn.Close(); err != nil {
		log.Printf("client failed to close connection: %s", err)
	}

	s.GracefulStop()
	if err := lis.Close(); err != nil {
		log.Printf("server listener failed to close: %s", err)
	}

	os.Exit(exitCode)
}

func TestRotateImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.RotateImage(ctx, &pb.NLImageRotateRequest{Rotation: 0, Image: &pb.NLImage{Color: true, Data: nil, Width: 0, Height: 0}})
	if err != nil {
		t.Fatalf("RotateImage failed: %v", err)
	}
}
