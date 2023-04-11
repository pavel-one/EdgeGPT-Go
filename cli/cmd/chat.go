package cmd

import (
	"bufio"
	"fmt"
	"github.com/MichaelMure/go-term-markdown"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/pavel-one/EdgeGPT-Go/internal/Logger"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	chat *EdgeGPT.GPT
	r    bool
)

var ChatCmd = &cobra.Command{
	Use:   "Chat",
	Short: "Edge Bing chat",
	Long:  "Simple cli for speaking with EdgeGPT Bing ",
	Run:   runChat,
}

func init() {
	rootCmd.AddCommand(ChatCmd)
	ChatCmd.Flags().BoolP("rich", "r", false, "parse markdown to terminal")
	logger = Logger.NewLogger("Chat")
	storage = EdgeGPT.NewStorage()
}

func runChat(cmd *cobra.Command, args []string) {
	rich, err := cmd.Flags().GetBool("rich")
	if err != nil {
		log.Fatalln(err)
	}
	r = rich

	newChat("chat")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nYou:\n    ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "q" || input == "quiet" {
			fmt.Println("Good bye!")
			break
		}

		ask(input)
	}
}

func newChat(key string) {
	gpt, err := storage.GetOrSet(key)
	if err != nil {
		logger.Fatalf("Failed to create new chat: %v", err)
	}
	chat = gpt
}

func ask(input string) {
	if r {
		rich(input)
		return
	} else {
		base(input)
		return
	}
}

func base(input string) {
	var l int

	mw, err := chat.AskAsync(input)
	if err != nil {
		logger.Fatalln(err)
	}

	go mw.Worker()

	for range mw.Chan {
		var res string
		ans := mw.Answer.GetAnswer()

		anslen := len(ans)

		if anslen == 0 {
			continue
		}

		if l == 0 {
			res = ans
		} else if 0 < l && l < anslen {
			res = ans[l:]
		}
		l = anslen
		fmt.Print(res)
	}
}

func rich(input string) {
	fmt.Println("Bot:")

	mw, err := chat.AskSync(input)
	if err != nil {
		logger.Fatalln(err)
	}

	ans := mw.Answer.GetAnswer()
	if ans == "" {
		return
	}

	result := markdown.Render(ans, 150, 4)

	if result == nil {
		fmt.Println(ans)
		return
	}

	fmt.Print(string(result))

	return
}
