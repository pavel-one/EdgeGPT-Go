package main

import (
	"EdgeGPT-Go/internal/EdgeGPT"
	"EdgeGPT-Go/internal/GRPC"
	pb "EdgeGPT-Go/pkg/GRPC/GPT"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
)

var storage = EdgeGPT.NewStorage()
var interceptor = func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	p, _ := peer.FromContext(ss.Context())
	if _, err := storage.GetOrSet(p.Addr.String()); err != nil {
		return err
	}

	err := handler(srv, ss)
	return err
}

func main() {
	srv := GRPC.NewServer(storage)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(interceptor),
	)
	pb.RegisterGptServiceServer(s, srv)
	log.Println("Starting server on port 8080")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
