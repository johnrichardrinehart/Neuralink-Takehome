package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
	"google.golang.org/grpc"
)

// server is used to implement image.NLImageServiceServer
type server struct {
	pb.UnimplementedNLImageServiceServer
}

func main() {
	var (
		port  string
		host  string
		debug bool
	)
	flag.StringVar(&port, "port", "2222", "server listening port (default: 2222)")
	flag.StringVar(&host, "host", "localhost", "default interface on which to listen (default: localhost)")
	flag.BoolVar(&debug, "debug", true, "debug mode (logs some runtime behavior)")
	flag.Parse()

	if v, err := strconv.Atoi(port); err != nil || v < 0 || v > 1<<16 {
		log.Fatal("invalid port number specified - must be between 0 and 65535 (0 requests a random port from the OS kernel)")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	if debug {
		log.Printf("gRPC server successfully created (listening at %s:%s)", host, port)
	}

	pb.RegisterNLImageServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
