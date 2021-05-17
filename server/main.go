package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	server "github.com/johnrichardrinehart/Neuralink-Takehome/server/pkg"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
	"google.golang.org/grpc"
)

func main() {
	var (
		port  string
		host  string
		debug bool
	)
	flag.StringVar(&port, "port", "2222", "server listening port")
	flag.StringVar(&host, "host", "localhost", "default interface on which to listen")
	flag.BoolVar(&debug, "debug", false, "debug mode (logs some runtime behavior)")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

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

	pb.RegisterNLImageServiceServer(s, &server.Server{
		Debug: debug,
	})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-done
	s.GracefulStop()
	if debug {
		log.Printf("server gracefully exiting")
	}
}
