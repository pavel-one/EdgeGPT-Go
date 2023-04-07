package cmd

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	gpt *EdgeGPT.GPT
	r   bool
)

var rootCmd = &cobra.Command{
	Use:   "EdgeGPT-Go",
	Short: "CLI for using edge bing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rich, err := cmd.Flags().GetBool("rich")
		if err != nil {
			log.Fatalln(err)
		}
		r = rich

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
		l          int
		code       bool
		shown      bool
		lexer      string
		codeLabel  bool
		codeSource string
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
		l = len(ans)

		if r {
			if res == "```" || res == "`\n" {
				code = true
				codeLabel = true
				shown = false
				continue
			}

			if codeLabel {
				codeLabel = false
				lexer = res
				fmt.Print(lexer, ":")
				fmt.Println()
				continue
			}

			if res == "``" {
				code = false
				continue
			}

			if res == "`\n\n" {
				continue
			}

			if code {
				codeSource += res
				continue
			}

			if code == false && codeSource != "" && shown == false && lexer != "" {
				if err := quick.Highlight(os.Stdout, codeSource, lexer, "terminal", "monokai"); err != nil {
					log.Fatalln(err)
				}
				fmt.Println()
				shown = true
				continue
			}
		}

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
	rootCmd.Flags().BoolP("rich", "r", false, "Colorize code if it exists")
}
