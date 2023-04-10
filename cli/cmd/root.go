package cmd

import (
	"fmt"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using edge bing",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func newChat(key string) *EdgeGPT.GPT {
	s := EdgeGPT.NewStorage()

	gpt, err := s.GetOrSet(key)
	if err != nil {
		log.Fatalln(err)
	}

	return gpt
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
