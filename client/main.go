package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"

	"google.golang.org/grpc"
)

func main() {
	var (
		port   string
		host   string
		input  string
		output string
		rotate string
		mean   bool
		// debug  bool
	)
	flag.StringVar(&port, "port", "2222", "port of server (default: 2222)")
	flag.StringVar(&host, "host", "localhost", "server host (default: localhost)")
	flag.StringVar(&input, "input", "in.jpg", "path to the input file (default: in.jpg)")
	flag.StringVar(&output, "output", "out.jpg", "path to the output file (default: out.jpg)")
	flag.StringVar(&rotate, "rotate", "NONE", "counterclockwise rotation angle: NONE, NINETY_DEG, ONE_EIGHTY_DEG, TWO_SEVENTY_DEG (default: NONE)")
	flag.BoolVar(&mean, "mean", false, "boolean option to apply mean filter to image")
	// flag.BoolVar(&debug, "debug", true, "debug mode (logs some runtime behavior)")
	flag.Parse()

	if v, err := strconv.Atoi(port); err != nil || v < 0 || v > 1<<16 {
		log.Fatal("invalid port number specified - must be between 0 and 65535 (0 requests a random port from the OS kernel)")
	}

	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("client failed to connect to server: %v", err)
	}

	defer conn.Close()

	_ = pb.NewNLImageServiceClient(conn)
}
