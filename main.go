package main

import (
	"context"
	"github.com/disintegration/gift"
	pb "github.com/vanbrabantf/microservice/ImageService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	log.Println("Call Main")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()

	pb.RegisterImageServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (s *server) GetImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
	p := updateImage()
	log.Println("Call")

	return &pb.ImageResponse{Path: p}, nil
}

func updateImage() string {
	src := loadImage("testdata/poster.png")

	filters := []gift.Filter{
		gift.Sigmoid(0.5, 7),
		gift.Pixelate(5),
		gift.Colorize(240, 50, 100),
		gift.Grayscale(),
		gift.Sepia(100),
		gift.Invert(),
		gift.Mean(5, true),
		gift.Median(5, true),
		gift.Minimum(5, true),
		gift.Maximum(5, true),
		gift.Hue(45),
		gift.ColorBalance(10, -10, -10),
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	g := gift.New(filters[r1.Intn(11)])
	dst := image.NewNRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	p := "testdata/updated.png"
	saveImage(p, dst)

	return p
}

func loadImage(filename string) image.Image {
	fImg1, _ := os.Open(filename)
	defer fImg1.Close()
	img1, _, _ := image.Decode(fImg1)

	return img1
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}
	log.Print("png.Encode worked!:", filename)
}
