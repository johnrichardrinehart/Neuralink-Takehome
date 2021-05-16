package server

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/BurntSushi/graphics-go/graphics"
	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
	"github.com/pixiv/go-libjpeg/rgb"
)

// server is used to implement image.NLImageServiceServer
type Server struct {
	pb.UnimplementedNLImageServiceServer
	Debug bool
}

// rgbDraw wraps an rgb.Image to implement the image.Drawer interface (used in rotation)
type rgbDraw struct {
	*rgb.Image
}

func (d *rgbDraw) Set(x, y int, c color.Color) {
	w := d.Rect.Max.X // width

	i := 3 * (y*w + x) // array index from (x,y)

	r, g, b, _ := c.RGBA() // alpha-scaled, need to reduce by 255

	d.Pix[i] = byte(r >> 8)
	d.Pix[i+1] = byte(g >> 8)
	d.Pix[i+2] = byte(b >> 8)
}

// RotateImage accepts a request to rotate an image and, if it's of a valid type (PNG/JPG/GIF)
// rotates it by the requested angle and returns it
// Failure to decode will return an error
func (s Server) RotateImage(ctx context.Context, req *pb.NLImageRotateRequest) (*pb.NLImage, error) {
	h := int(req.Image.Height) // WARNING: int32 -> int; should be fine for most systems
	w := int(req.Image.Width)  // WARNING: int32 -> int; should be fine for most systems
	c := req.Image.Color

	if s.Debug {
		log.Printf("received request to rotate a %dx%d image: %v degrees", h, w, 90*req.Rotation)
	}

	if err := validateImage(req.Image); err != nil {
		return nil, fmt.Errorf("image failed to validate: %s", err)
	}

	if req.Rotation == 0 {
		return req.Image, nil
	}

	box := image.Rect(0, 0, w, h)
	if !c {
		src := image.NewGray(box)
		src.Pix = req.Image.Data
		dst := image.NewGray(box)

		graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: -1 * (math.Pi / 2) * float64(req.Rotation)})

		req.Image.Data = dst.Pix
	} else {
		src := rgb.NewImage(box)
		src.Pix = req.Image.Data
		dst := &rgbDraw{rgb.NewImage(box)}

		graphics.Rotate(dst, src, &graphics.RotateOptions{Angle: -1 * (math.Pi / 2) * float64(req.Rotation)})

		req.Image.Data = dst.Pix
	}

	return req.Image, nil
}

func (s Server) MeanFilter(ctx context.Context, img *pb.NLImage) (*pb.NLImage, error) {
	if s.Debug {
		log.Printf("received request to filter an image")
	}

	if err := validateImage(img); err != nil {
		return nil, fmt.Errorf("image failed to validate: %s", err)
	}

	stride := 3
	if !img.Color {
		stride = 1
	}

	var neighbors [][2]int = [][2]int{
		{-1 * stride, 0},           // bottom
		{1 * stride, 0},            // top
		{0, -1 * stride},           // left
		{0, 1 * stride},            // right
		{-1 * stride, -1 * stride}, // bottom left
		{1 * stride, 1 * stride},   // top right
		{-1 * stride, 1 * stride},  // top left
		{1 * stride, -1 * stride},  // bottom right
	}

	for i := 0; i < len(img.Data); i += 1 {
		row := i / int(img.Width)
		col := i % int(img.Width)
		acc := img.Data[i]
		var cnt int
		for _, n := range neighbors {
			ro := n[0]                                // row offset
			co := n[1]                                // col offset
			j := (row-ro)*int(img.Width) + (col + co) // WARNING: cast could break on extremely large widths
			if j < 0 || j > len(img.Data)-1 {
				continue
			}
			cnt += 1
			acc += img.Data[j]
		}
		img.Data[i] = byte(int(acc) / cnt) // mean
	}

	return img, nil
}

func validateImage(img *pb.NLImage) error {
	h := int(img.Height) // WARNING: int32 -> int; should be fine for most systems
	w := int(img.Width)  // WARNING: int32 -> int; should be fine for most systems
	c := img.Color

	if (c && (len(img.Data) != 3*h*w)) || (!c && (len(img.Data) != h*w)) {
		return errors.New("invalid data length - should be a 3x or 1x multiple of height*width")
	}

	return nil
}
