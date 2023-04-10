package cmd

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/pavel-one/EdgeGPT-Go"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	gpt *EdgeGPT.GPT
	r   bool
)

var EdgeGPTCliCmd = &cobra.Command{
	Use:   "EdgeGPTCli",
	Short: "CLI chat for using edge bing",
	Long:  ``,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(EdgeGPTCliCmd)
	EdgeGPTCliCmd.Flags().BoolP("rich", "r", false, "Colorize code if it exists")
}

func run(cmd *cobra.Command, args []string) {
	rich, err := cmd.Flags().GetBool("rich")
	if err != nil {
		log.Fatalln(err)
	}
	r = rich

	gpt = newChat("cli")
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Hello, I am a chatbot for speaking with edge bing!")

	for {
		fmt.Print("\nYou:\n    ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "q" || input == "quiet" {
			fmt.Println("Goodbye!")
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

	regex := regexp.MustCompile("([\\s\\S]*?)```([a-zA-Z]+[+]*)([\\s\\S]*?)```([^`]+)")
	matches := regex.FindAllStringSubmatch(ans, -1)

	if matches == nil {
		fmt.Print(ans)
		return
	}
	if matches[0] == nil {
		fmt.Print(ans)
		return
	}

	for _, m := range matches {
		fmt.Print(m[1])

		if err := quick.Highlight(os.Stdout, m[3], m[2], "terminal", "monokai"); err != nil {
			log.Fatalln(err)
		}

		fmt.Print(m[4])
	}

	return
}
