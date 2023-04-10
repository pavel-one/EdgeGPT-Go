package cmd

import (
	"bufio"
	"fmt"
	"github.com/MichaelMure/go-term-markdown"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	gpt *EdgeGPT.GPT
	r   bool
)

var EdgeGPTCliCmd = &cobra.Command{
	Use:   "Chat",
	Short: "EdgeBing chat",
	Long:  "Simple cli for speaking with EdgeGPT Bing ",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(EdgeGPTCliCmd)
	EdgeGPTCliCmd.Flags().BoolP("rich", "r", false, "parse markdown to terminal")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go handleSignal(sigs)
}

func run(cmd *cobra.Command, args []string) {
	rich, err := cmd.Flags().GetBool("rich")
	if err != nil {
		log.Fatalln(err)
	}
	r = rich

	gpt = newChat("cli")
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

	mw, err := gpt.AskAsync(input)
	if err != nil {
		log.Fatalln(err)
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

	mw, err := gpt.AskSync(input)
	if err != nil {
		log.Fatalln(err)
	}

	ans := mw.Answer.GetAnswer()
	if ans == "" {
		return
	}

	result := markdown.Render(ans, 100, 1)

	if result == nil {
		fmt.Println(ans)
		return
	}

	fmt.Println(string(result))

	return
}
