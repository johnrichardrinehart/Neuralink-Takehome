package server

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
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

func TestMeanFilter(t *testing.T) {
	tt := []struct {
		name         string
		img          *pb.NLImage
		expBytes     []byte
		expError     bool
		errSubstring string
	}{
		{
			"3x3 - no color",
			&pb.NLImage{
				Color:  false,
				Data:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8},
				Width:  3,
				Height: 3,
			},
			[]byte{
				(0 + 1 + 3 + 4) / 4,
				(0 + 1 + 2 + 3 + 4 + 5) / 6,
				(1 + 2 + 4 + 5) / 4,
				(0 + 1 + 3 + 4 + 6 + 7) / 6,
				(0 + 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8) / 9,
				(1 + 2 + 4 + 5 + 7 + 8) / 6,
				(3 + 4 + 6 + 7) / 4,
				(3 + 4 + 5 + 6 + 7 + 8) / 6,
				(4 + 5 + 7 + 8) / 4,
			},
			false,
			"",
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.MeanFilter(ctx, tst.img)

			if err != nil {
				if !tst.expError {
					t.Fatalf("encountered an error when none was expected: %s", err)
				}
				if !strings.Contains(err.Error(), tst.errSubstring) {
					t.Fatalf("unexpected error string encountered - expected %s - got %s", tst.errSubstring, err)
				}
			} else {
				if tst.expError {
					t.Fatalf("failed to encounter an error when the following was expected: %s", tst.errSubstring)
				}
			}

			for i, v := range resp.Data {
				if v != resp.Data[i] {
					t.Fatal("failed to calculate the mean properly")
				}
			}
		})
	}
}
