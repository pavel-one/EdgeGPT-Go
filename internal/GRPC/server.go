package GRPC

import (
	"EdgeGPT-Go/internal/EdgeGPT"
	pb "EdgeGPT-Go/pkg/GRPC/GPT"
	"google.golang.org/grpc/peer"
	"log"
	"time"
)

type Server struct {
	pb.UnimplementedGptServiceServer
	Storage *EdgeGPT.Storage
}

func NewServer(s *EdgeGPT.Storage) *Server {
	return &Server{
		Storage: s,
	}
}

func (s *Server) Ask(r *pb.AskRequest, stream pb.GptService_AskServer) error {
	p, _ := peer.FromContext(stream.Context())
	gpt, err := s.Storage.GetOrSet(p.Addr.String())
	if err != nil {
		return err
	}

	message, err := gpt.AskAsync(r.Text)
	if err != nil {
		return err
	}

	go func() {
		err := message.Worker()
		if err != nil {
			log.Println(err)
		}
	}()

	time.Sleep(time.Second)

	for _ = range message.Chan {
		msg := message.Answer.GetAnswer()
		if msg == "" {
			continue
		}

		res := &pb.AskResponse{
			Text:       message.Answer.GetAnswer(),
			MaxUnit:    uint64(message.Answer.GetMaxUnit()),
			UnitUser:   uint64(message.Answer.GetUserUnit()),
			ExpiryTime: uint64(gpt.ExpiredAt.Unix()),
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}
