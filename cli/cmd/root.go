package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var gpt *EdgeGPT.GPT

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using edge bing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gpt = newChat()
		reader := bufio.NewReader(os.Stdin)

		color.Green("Hello, I am a chatbot for speak with edge bing!")

		for {
			fmt.Print("You:\n    ")
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
	var (
		l int
		//code bool
		//labelCount int
		//codeSource string
	)
	mw, err := gpt.AskAsync(input)
	if err != nil {
		log.Fatalln(err)
	}

	go mw.Worker()

	fmt.Print("Bot:\n    ")
	for range mw.Chan {
		ans := mw.Answer.GetAnswer()

		if len(ans) == 0 {
			continue
		}

		res := ans[l:]
		if res == "```" || res == "go" || res == " \n```\n\n" {
			//code = true
			//l = len(ans)
			//continue
		}

		//if res == ""
		//if res == "```" {
		//fmt.Println("```HERE")
		//code = !code
		//l = len(ans)
		//continue
		//}

		//if code {
		//	codeSource += res
		//	l = len(ans)
		//	continue
		//}
		//
		//if code == false && codeSource != "" {
		//	if err := quick.Highlight(os.Stdout, codeSource, "go", "terminal", "monokai"); err != nil {
		//		log.Fatalln(err)
		//	}
		//	l = len(ans)
		//	continue
		//}

		l = len(ans)
		fmt.Print(res)
	}
}

func newChat() *EdgeGPT.GPT {
	s := EdgeGPT.NewStorage()

	gpt, err := s.GetOrSet("cli")
	if err != nil {
		log.Fatalln(err)
	}

	return gpt
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
