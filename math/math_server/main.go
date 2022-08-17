package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/AlexanderGuan/GRPC-TEST/math/mathpb"
	"google.golang.org/grpc"
)

type server struct {
	*pb.UnimplementedMathServer
}

func (s *server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Printf("--- gRPC Unary RPC ---\n")
	fmt.Printf("request received: %v\n", in)
	return &pb.SumResponse{Result: in.FirstNum + in.SecondNum}, nil
}

func (s *server) PrimeFactors(in *pb.PrimeFactorsRequest, stream pb.Math_PrimeFactorsServer) error {
	fmt.Printf("--- gRPC Server-side Streaming RPC ---\n")
	fmt.Printf("request received: %v\n", in)

	num := in.Num
	factor := int64(2)

	for num > 1 {
		if num%factor == 0 {
			stream.Send(&pb.PrimeFactorsResponse{Result: factor})
			num = num / factor
			continue
		}
		factor++
	}
	return nil
}

func main() {
	port := flag.Int("port", 50051, "the port to serve on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterMathServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
