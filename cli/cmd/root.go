package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/pavel-one/EdgeGPT-Go/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using edge bing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		color.Green("Hello, I am a chatbot for speak with edge bing!")

		for {
			fmt.Println("You:")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" || input == "q" || input == "quiet" {
				color.Yellow("Goodbye!")
				break
			}

			ask(input)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ask(input string) {
	gpt := gpt()

	mw, err := gpt.AskSync(input)
	if err != nil {
		log.Fatalln(err)
	}

	for range mw.Chan {
	}
	color.Cyan(mw.Answer.GetAnswer())
}

func gpt() *EdgeGPT.GPT {
	conf, err := config.NewGpt()
	if err != nil {
		log.Fatalln(err)
	}

	gpt, err := EdgeGPT.NewGPT(conf)
	if err != nil {
		log.Fatalln(err)
	}

	return gpt
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
