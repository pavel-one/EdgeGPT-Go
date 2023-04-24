package cmd

import (
	"bufio"
	"fmt"
	term_markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	chat            *EdgeGPT.GPT
	r               bool
	toHtml          bool
	output          string
	withoutTerminal bool
	style           string
)

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Edge Bing chat",
	Long:  "Simple cli for speaking with EdgeGPT Bing ",
	Run:   runChat,
}

func init() {
	rootCmd.AddCommand(ChatCmd)
	ChatCmd.Flags().BoolVarP(&r, "rich", "r", false, "parse markdown to terminal")
	ChatCmd.Flags().BoolVarP(&toHtml, "html", "", false, "parse markdown to html(use with --output)")
	ChatCmd.Flags().StringVarP(&output, "output", "o", "", "output file(markdown or html like test.md or test.html or just text like `test` file)")
	ChatCmd.Flags().BoolVarP(&withoutTerminal, "without-term", "w", false, "if output set will be write response to file without terminal")
	ChatCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "set endpoint for create conversation(if the default one doesn't suit you)")
	ChatCmd.Flags().StringVarP(&style, "style", "s", "balanced", "set conversation style(creative, balanced, precise)")
}

func runChat(cmd *cobra.Command, args []string) {
	initLoggerWithStorage("Chat")
	setConversationEndpoint()
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
	if output != "" && withoutTerminal {
		ans := getAnswer(input)
		writeWithFlags([]byte(ans))
		return
	}

	if r {
		rich(input)
	} else {
		base(input)
	}
}

func base(input string) {
	var l int

	mw, err := chat.AskAsync(style, input)
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

	go writeWithFlags([]byte(mw.Answer.GetAnswer()))

	return
}

func rich(input string) {
	fmt.Println("Bot:")

	ans := getAnswer(input)

	go writeWithFlags([]byte(ans))

	result := term_markdown.Render(ans, 150, 4)

	if result == nil {
		fmt.Println(ans)
		return
	}

	fmt.Print(string(result))

	return
}

func getAnswer(input string) string {
	mw, err := chat.AskSync(style, input)
	if err != nil {
		logger.Fatalln(err)
	}

	ans := mw.Answer.GetAnswer()
	if ans == "" {
		logger.Fatalf("failed to get answer")
	}

	return ans
}

func writeWithFlags(data []byte) {
	if output != "" {
		if toHtml {
			d := mdToHtml(data)
			writeToFile(d)
		} else {
			writeToFile(data)
		}
	}
}

func mdToHtml(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func writeToFile(data []byte) {
	_, err := os.Stat(output)
	if os.IsNotExist(err) {
		dir := filepath.Dir(output)
		if dir != "." {
			if err = os.MkdirAll(dir, 0755); err != nil {
				logger.Fatalf("failed to create dir %s: %v", dir, err)
			}
		}

		if err := os.WriteFile(output, data, 0644); err != nil {
			logger.Fatalf("failed while write data to file `%s`: %v", output, err)
		}
	} else {
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0755)
		defer file.Close()

		if err != nil {
			logger.Fatalf("failed to open file `%s`: %v", output, err)
		}

		_, err = file.Write(data)
		if err != nil {
			logger.Fatalf("failed to write string to file `%s`: %v", file.Name(), err)
		}

		if err = file.Sync(); err != nil {
			logger.Fatalf("failed to sync data for file `%s`: %v", file.Name(), err)
		}
	}
	logger.Info("Response written successfully to " + output)
}
