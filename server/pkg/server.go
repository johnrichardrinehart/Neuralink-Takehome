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
	w := int(img.Width)
	h := int(img.Height)

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
		{-1, -1}, // top left
		{0, -1},  // top
		{1, -1},  // top right
		{-1, 0},  // left
		{1, 0},   // right
		{-1, 1},  // bottom left
		{0, 1},   // bottom
		{1, 1},   // bottom right
	}

	mean := make([]byte, len(img.Data))

	// pixel loop
	for i := 0; i < len(img.Data); i += stride {
		ix, iy := iToXY(i, stride, w) // position of 'r' in 'rgb'

		// color loop
		for c := 0; c < stride; c += 1 {
			acc := int(img.Data[i+c]) // accumulator
			cnt := 1                  // initial weight
			for _, n := range neighbors {
				co := n[0] // col displacement
				ro := n[1] // row displacement

				// neighbor index
				if ix+co < 0 || ix+co > w-1 {
					continue
				}
				if iy+ro < 0 || iy+ro > h-1 {
					continue
				}

				ni := XYToI(ix+co, iy+ro, stride, c, w)

				cnt += 1
				acc += int(img.Data[ni])
			}
			// calculate the mean
			mean[i+c] = byte(int(acc) / cnt)
		}

	}

	img.Data = mean

	return img, nil
}

// iToXY takes the index of an element in a row-ordered linear array and outputs the (x,y) coordinate of the pixel
// within the picture
// each pixel is assumed to be a stride-dimensional sub-array
func iToXY(i, stride, width int) (x, y int) {
	n := i / stride // aggregate by stride length (rgb have the same pixel)
	x = n % width
	y = n / width
	return
}

// XYtoI converts the (x,y) coordinate of a stride-dimensional pixel and a color "offset"
// and returns the index in a row-ordered linear array
func XYToI(x, y, stride, colorOffset, width int) (i int) {
	i = stride*(x+y*width) + colorOffset
	return
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
