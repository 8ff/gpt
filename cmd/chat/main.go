package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	gpt "github.com/8ff/gpt/pkg/gpt_3_5_turbo"
)

func main() {
	// Read API_TOKEN from env
	token := os.Getenv("API_TOKEN")

	api, err := gpt.Init(gpt.Params{
		API_TOKEN:          token,
		KeepMessageHistory: true,
		StripNewline:       true,
	})
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a message (type 'exit' to quit): ")
		scanner.Scan()
		input := scanner.Text()

		fmt.Printf("USER: %s\n", input)

		if strings.ToLower(input) == "exit" {
			break
		}

		choices, err := api.Query(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, choice := range choices {
			fmt.Printf("GPT: %s\n", choice.Message.Content)
		}

		fmt.Println("********************************************")
	}
}
