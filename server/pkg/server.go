package server

import (
	"context"
	"log"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
)

// server is used to implement image.NLImageServiceServer
type Server struct {
	pb.UnimplementedNLImageServiceServer
	Debug bool
}

func (s Server) RotateImage(ctx context.Context, req *pb.NLImageRotateRequest) (*pb.NLImage, error) {
	if s.Debug {
		log.Printf("received request to rotate an image: %v degrees", 90*req.Rotation)
	}

	return req.Image, nil
}

func (s Server) MeanFilter(ctx context.Context, img *pb.NLImage) (*pb.NLImage, error) {
	if s.Debug {
		log.Printf("received request to filter an image")
	}

	return img, nil
}
