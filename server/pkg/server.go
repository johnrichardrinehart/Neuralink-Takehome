package server

import (
	"context"
	"errors"
	"fmt"
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

	h := int(req.Image.Height) // WARNING: int32 -> int; should be fine for most systems
	w := int(req.Image.Width)  // WARNING: int32 -> int; should be fine for most systems
	c := req.Image.Color

	if err := validateImage(req.Image); err != nil {
		return nil, fmt.Errorf("image failed to validate: %s", err)
	}

	// coerce the input image into RGBa format to re-use stdlib
	var imga []byte

	if !c {
		imga = make([]byte, 2*h*w)
		for i, v := range req.Image.Data {
			imga[i] = v // optimistic
			imga[i+1] = 1 << 2
		}
	} else {
		imga = make([]byte, 4*h*w/3)
		for i, v := range req.Image.Data {
			imga[i] = v // optimistic
			if i%4 == 0 {
				imga[i] = 1 << 3
			}
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, int(req.Image.Width), int(req.Image.Height)))
	img.Pix = imga

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

	if err := validateImage(img); err != nil {
		return nil, fmt.Errorf("image failed to validate: %s", err)
	}

	return img, nil
}

func validateImage(img *pb.NLImage) error {
	h := int(img.Height) // WARNING: int32 -> int; should be fine for most systems
	w := int(img.Width)  // WARNING: int32 -> int; should be fine for most systems
	c := img.Color

	if (c && len(img.Data) != 3*h*w) || len(img.Data) != h*w {
		return errors.New("invalid data length - should be a 3x or 1x multiple of height*width")
	}

	return nil
}
