package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"image/jpeg"
	"image/png"

	pb "github.com/johnrichardrinehart/Neuralink-Takehome/proto"
	"github.com/pixiv/go-libjpeg/rgb"

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
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure(), grpc.WithBlock())
	cancel() // avoid resource leak
	if err != nil {
		log.Fatalf("client failed to connect to server: %v", err)
	}

	defer conn.Close()

	fin, err := os.Open(input)
	if err != nil {
		log.Fatalf("failed to open file %s: %s", input, err)
	}

	img, ft, err := image.Decode(fin)
	if err != nil {
		log.Fatalf("failed to decode input file: %s", err)
	}

	log.Printf("discovered filetype: %s", ft)

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	bs := make([]byte, w*h*3)
	for j := 0; j < h-1; j++ {
		for i := 0; i < w-1; i++ {
			color := img.At(i, j)
			r, g, b, _ := color.RGBA()
			n := 3 * (i + j*w)

			bs[n] = byte(r >> 8)
			bs[n+1] = byte(g >> 8)
			bs[n+2] = byte(b >> 8)
		}
	}

	nlImg := pb.NLImage{
		Color:  true,
		Data:   bs,
		Width:  int32(w),
		Height: int32(h),
	}

	c := pb.NewNLImageServiceClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if mean {
		resp, err := c.MeanFilter(ctx, &nlImg)
		if err != nil {
			log.Fatalf("failed to apply mean filter: %s", err)
		}

		nlImg = pb.NLImage{
			Color:  resp.Color,
			Data:   resp.Data,
			Width:  resp.Width,
			Height: resp.Height,
		}
	}

	var rotation pb.NLImageRotateRequest_Rotation
	switch rotate {
	case "NONE":
		rotation = pb.NLImageRotateRequest_NONE
	case "NINETY_DEG":
		rotation = pb.NLImageRotateRequest_NINETY_DEG
	case "ONE_EIGHTY_DEG":
		rotation = pb.NLImageRotateRequest_ONE_EIGHTY_DEG
	case "TWO_SEVENTY_DEG":
		rotation = pb.NLImageRotateRequest_TWO_SEVENTY_DEG
	}

	req := pb.NLImageRotateRequest{
		Rotation: rotation,
		Image:    &nlImg,
	}

	resp, err := c.RotateImage(ctx, &req)
	if err != nil {
		log.Printf("failed to rotate image: %s", err)
	}

	fout, err := os.Create(output)
	if err != nil {
		log.Fatalf("failed to create output file %s: %s", output, err)
	}

	m := rgb.NewImage(image.Rect(0, 0, w, h))
	m.Pix = resp.Data

	switch ft {
	case "jpeg":
		if err := jpeg.Encode(fout, m, &jpeg.Options{Quality: 100}); err != nil {
			log.Fatalf("failed to write jpeg to file %s: %s", output, err)
		}
	case "png":
		if err := png.Encode(fout, m); err != nil {
			log.Fatalf("failed to write jpeg to file %s: %s", output, err)
		}
	}
}
