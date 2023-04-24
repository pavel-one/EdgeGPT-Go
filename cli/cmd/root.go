package cmd

import (
	"fmt"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/pavel-one/EdgeGPT-Go/internal/Logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger  *zap.SugaredLogger
	storage *EdgeGPT.Storage
)

var endpoint string

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using Edge Bing",
	Long:  "Cli for using Edge Bing. Available commands:\nChat - for speaking with Bing\ngRPC - start gRPC server for speaking with Bing",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go handleSignal(sigs)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setConversationEndpoint() {
	if endpoint != "" {
		if err := os.Setenv("CONVERSATION_ENDPOINT", endpoint); err != nil {
			logger.Fatal("failed to set CONVERSATION_ENDPOINT to environment")
		}
	}
}

func initLoggerWithStorage(channel string) {
	logger = Logger.NewLogger(channel)
	storage = EdgeGPT.NewStorage()
}

func handleSignal(sigs chan os.Signal) {
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGINT:
			fmt.Println("\nGood bye!")
			os.Exit(0)
		}
	}
}
