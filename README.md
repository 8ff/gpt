![logo](media/logo.svg "GPT-3.5 Turbo & GPT-4 Chatbot Golang Library")
# GPT-3.5 Turbo & GPT-4 Chatbot Golang Library
[![Go Report Card](https://goreportcard.com/badge/github.com/8ff/gpt)](https://goreportcard.com/report/github.com/8ff/gpt)
[![GoDoc](https://godoc.org/github.com/8ff/gpt?status.svg)](https://godoc.org/github.com/8ff/gpt)
[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://github.com/8ff/gpt/blob/main/LICENSE)

This repository contains a Golang implementation of a chatbot using the OpenAI GPT-3.5 Turbo as well as the new GPT-4 API. The chatbot is able to generate human-like responses to user queries.

## Library
The core of this chatbot implementation are the gpt_3_5_turbo and gpt_4 packages, which are Golang libraries for interacting with the OpenAI GPT-3.5 Turbo and GPT-4 API. These libraries provide a simple API for sending text queries to the GPT-3.5 and GPT-4 models and receiving human-like responses in return. They include support for features like setting the API token, configuring the request, and managing message history. With these libraries, developers can easily incorporate the power of the GPT-3.5 & GPT-4 models into their Golang applications and build intelligent chatbots or other NLP-driven tools.


## Prerequisites

To use this chatbot, you need to have an API token for the OpenAI GPT-3.5 Turbo or GPT-4 API. You can obtain one by following the instructions on the [OpenAI website](https://beta.openai.com/signup/).

## Example
```go
package main

import (
	"fmt"
	"os"

	gpt "github.com/8ff/gpt/pkg/gpt_3_5_turbo"
	// gpt "github.com/8ff/gpt/pkg/gpt_4"
)

func main() {
	// Read API_TOKEN from env
	token := os.Getenv("API_TOKEN")

	api, err := gpt.Init(gpt.Params{
		API_TOKEN:    token,
		StripNewline: true,
		Request: gpt.ChatRequest{
			Model: "gpt-3.5-turbo",
			// Model: "gpt-4-0613",
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
```

## Installation

To install this chatbot, you need to have Go installed on your machine. Once you have Go installed, you can use the following command to download and install the chatbot:

```bash
git clone github.com/8ff/gpt
```

## Usage

This repository also contains two demo applications for using the chatbot: a single request app and a chat app.

### Single Request App

The single request app is located in `cmd/singleRequest`. To use the single request app, set the `API_TOKEN` environment variable to your OpenAI API token, and run the following command:

```bash
go run main.go
```

![](media/singleRequest.gif "Single Request App")


The app will prompt you for a message, and generate a response based on your input.

### Chat App

The chat app is located in `cmd/chat`. To use the chat app, set the `API_TOKEN` environment variable to your OpenAI API token, and run the following command:

```bash
go run main.go
```

![](media/chat.gif "Chat App")


The app will prompt you for a message, and generate a response based on your input. You can continue chatting with the bot until you type "exit".

## License

This code is released under the GPL3 License. See `LICENSE` for more information.
