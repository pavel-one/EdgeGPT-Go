package cmd

import (
	"github.com/pavel-one/EdgeGPT-Go/internal/GRPC"
	pb "github.com/pavel-one/EdgeGPT-Go/pkg/GRPC/GPT"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

var gRPCCmd = &cobra.Command{
	Use:   "gRPC",
	Short: "Start gRPC server",
	Long:  `Command for starting gRPC server for speaking with Edge Bing`,
	Run:   rungRPC,
}

func init() {
	rootCmd.AddCommand(gRPCCmd)
	gRPCCmd.Flags().StringP("port", "p", "8080", "port for gRPC server")
}

func rungRPC(cmd *cobra.Command, args []string) {
	initLoggerWithStorage("gRPC")
	port, err := cmd.Flags().GetString("port")
	if err != nil {
		logger.Fatalf("failed to get flag `port`: %v", err)
	}

	srv := GRPC.NewServer(storage)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGptServiceServer(s, srv)
	logger.Infoln("Starting server on port " + port)

	if err := s.Serve(listener); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
