package main

import (
	"fmt"
	"os"

	gpt "github.com/8ff/gpt/pkg/gpt_3_5_turbo"
)

func main() {
	// Read API_TOKEN from env
	token := os.Getenv("API_TOKEN")

	api, err := gpt.Init(gpt.Params{
		API_TOKEN:    token,
		StripNewline: true,
		Request: gpt.ChatRequest{
			Model: "gpt-3.5-turbo",
		},
	})
	if err != nil {
		panic(err)
	}

	choices, err := api.Query("What are you ?")
	if err != nil {
		panic(err)
	}

	for _, choice := range choices {
		fmt.Printf("Response: %s\n", choice.Message.Content)
	}
}
