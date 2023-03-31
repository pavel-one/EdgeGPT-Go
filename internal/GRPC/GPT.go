package GRPC

import (
	pb "EdgeGPT-Go/pkg/GRPC/GPT"
	"google.golang.org/grpc/peer"
	"log"
)

type GPTServer struct {
	pb.UnimplementedGptServiceServer
}

func NewGPTServer() *GPTServer {
	return &GPTServer{}
}

func (s *GPTServer) Ask(r *pb.Empty, ss pb.GptService_AskServer) error {
	p, _ := peer.FromContext(ss.Context())
	log.Printf("Into method %s", p.Addr)
	return nil
}
