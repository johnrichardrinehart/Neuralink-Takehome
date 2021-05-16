package server

import (
	"bytes"
	"context"
	"errors"
	"image"
	"log"
	"math"

	"github.com/BurntSushi/graphics-go/graphics"
	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
)

// server is used to implement image.NLImageServiceServer
type Server struct {
	pb.UnimplementedNLImageServiceServer
	Debug bool
}

// RotateImage accepts a request to rotate an image and, if it's of a valid type (PNG/JPG/GIF)
// rotates it by the requested angle and returns it
// Failure to decode will return an error
func (s Server) RotateImage(ctx context.Context, req *pb.NLImageRotateRequest) (*pb.NLImage, error) {
	if s.Debug {
		log.Printf("received request to rotate an image: %v degrees", 90*req.Rotation)
	}

	rdr := bytes.NewBuffer(req.Image.Data)
	img, _, err := image.Decode(rdr)

	if err != nil {
		return nil, errors.New("failed to decode the image - invalid format")
	}

	if req.Rotation == 0 {
		return req.Image, nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(req.Image.Width), int(req.Image.Height)))
	graphics.Rotate(dst, img, &graphics.RotateOptions{Angle: (math.Pi / 2) * float64(req.Rotation)})

	req.Image.Data = dst.Pix

	return req.Image, nil
}

func (s Server) MeanFilter(ctx context.Context, img *pb.NLImage) (*pb.NLImage, error) {
	if s.Debug {
		log.Printf("received request to filter an image")
	}

	return img, nil
}
