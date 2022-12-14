// Package main implements a server for Echo service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/AlexanderGuan/GRPC-TEST/features/echopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server is used to implements echopb.EchoServer.
type server struct {
	*pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("--- gRPC Unary RPC ---\n")
	fmt.Printf("requetst received: %v\n", req)

	// Suppose it takes a long time to process(3 seconds)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled { // 表示客户端制动取消调用此RPC方法了，比如客户端按ctrl+c取消运行
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Errorf(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	return &pb.EchoResponse{Message: req.GetMessage()}, nil
}

func (s *server) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStreamingEcho not implemented")
}
func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStreamingEcho not implemented")
}
func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	return status.Errorf(codes.Unimplemented, "method BidirectionalStreamingEcho not implemented")
}

func main() {
	port := flag.Int("port", 50051, "the port to serve on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port)) // Specify the port we want to use to listen for client requests
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()                // Create an instance of the gRPC server
	pb.RegisterEchoServer(s, &server{})  // Register our service implementation with the gRPC server
	if err := s.Serve(lis); err != nil { // Call Serve() on the server with our port details to do a blocking wait until the process is killed or Stop() is called.
		log.Fatalf("failed to serve: %v", err)
	}
}
