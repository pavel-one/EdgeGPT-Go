package main

import (
	"EdgeGPT-Go/internal/GRPC"
	pb "EdgeGPT-Go/pkg/GRPC/GPT"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"time"
)

func main() {
	srv := GRPC.NewServer()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			// Обработка инициализации фактического подключения клиента
			p, _ := peer.FromContext(ss.Context())
			log.Println("peer address:", p.Addr)

			log.Printf("Client connected")
			err := handler(srv, ss)
			log.Printf("Client disconnected")
			return err
		}),
	)
	pb.RegisterGptServiceServer(s, srv)
	log.Println("Starting server...")

	go func() {
		time.Sleep(time.Second * 5)
		for i := 0; i < 5; i++ {
			conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
			if err != nil {
				panic(err)
			}

			client := pb.NewGptServiceClient(conn)
			for i := 0; i < 5; i++ {
				_, err := client.Ask(context.Background(), &pb.Empty{})
				if err != nil {
					log.Fatalf("Ошибка при вызове метода Ask: %v", err)
				}

			}
			conn.Close()
		}
	}()

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
