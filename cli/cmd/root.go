package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using edge bing",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
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
