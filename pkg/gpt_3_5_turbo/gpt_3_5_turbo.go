package gpt_3_5_turbo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Params struct {
	API_TOKEN          string      `json:"api_token,omitempty"`            // This will be stripped from the request before sending to the API. Required.
	StripNewline       bool        `json:"strip_newline,omitempty"`        // If set, the API will strip newlines from the beginning of the generated text. Optional, defaults to false.
	Request            ChatRequest `json:"request,omitempty"`              // The request body to send to the API. Required.
	KeepMessageHistory bool        `json:"keep_message_history,omitempty"` // If set, the message history will be kept in the response. Optional, defaults to false.
	MessageHistory     []Message   `json:"message_history,omitempty"`      // The message history to use. Optional, defaults to null.
}

type ChatRequest struct {
	Model            string             `json:"model"`                       // ID of the model to use. Currently, only gpt-3.5-turbo and gpt-3.5-turbo-0301 are supported. Required.
	Messages         []Message          `json:"messages"`                    // The messages to generate chat completions for, in the chat format. Required.
	Temperature      float64            `json:"temperature,omitempty"`       // What sampling temperature to use, between 0 and 2. Optional, defaults to 1.
	TopP             float64            `json:"top_p,omitempty"`             // An alternative to sampling with temperature, where the model considers the results of the tokens with top_p probability mass. Optional, defaults to 1.
	N                int                `json:"n,omitempty"`                 // How many chat completion choices to generate for each input message. Optional, defaults to 1.
	Stream           bool               `json:"stream,omitempty"`            // *** NOT IMPLEMENTED *** If set, partial message deltas will be sent. Optional, defaults to false.
	Stop             interface{}        `json:"stop,omitempty"`              // Up to 4 sequences where the API will stop generating further tokens. Optional, defaults to null.
	MaxTokens        int                `json:"max_tokens,omitempty"`        // The maximum number of tokens allowed for the generated answer. Optional, defaults to inf.
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`  // Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far. Optional, defaults to 0.
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"` // Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far. Optional, defaults to 0.
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`        // Modify the likelihood of specified tokens appearing in the completion. Optional, defaults to null.
	User             string             `json:"user,omitempty"`              // A unique identifier representing your end-user. Optional.
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (params *Params) Query(msg string) ([]Choice, error) {
	// If history is kept, copy all messages to request
	if params.KeepMessageHistory {
		params.Request.Messages = params.MessageHistory
	}

	// Append the message to the request
	params.Request.Messages = append(params.Request.Messages, Message{
		Role:    "user",
		Content: msg,
	})

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(params.Request)
	if err != nil {
		return nil, err
	}

	// Define the request object
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+params.API_TOKEN)

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If response is not 200 OK, return an error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OpenAI API returned status code %d", resp.StatusCode)
	}

	// Decode the response body into a struct
	var responseStruct ChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&responseStruct)
	if err != nil {
		return nil, err
	}

	// Append response to message
	for i, choice := range responseStruct.Choices {
		// If StripNewline is true, remove the 2 consecutive newlines
		if params.StripNewline {
			if len(choice.Message.Content) < 2 {
				continue
			}

			// If the first 2 characters are newlines, remove them
			if choice.Message.Content[:2] == "\n\n" {
				choice.Message.Content = choice.Message.Content[2:]
			}
		}

		// Store modified message back in response
		responseStruct.Choices[i].Message = choice.Message
	}

	// If history is kept, append the response to the message history
	if params.KeepMessageHistory {
		params.MessageHistory = append(params.MessageHistory, Message{
			Role:    "user",
			Content: msg,
		})

		for _, choice := range responseStruct.Choices {
			params.MessageHistory = append(params.MessageHistory, choice.Message)
		}
	}

	return responseStruct.Choices, nil
}

func Init(userParams Params) (*Params, error) {
	params := &Params{}

	// Check if API token is set
	if userParams.API_TOKEN == "" {
		return nil, fmt.Errorf("API_TOKEN is not set")
	} else {
		params.API_TOKEN = userParams.API_TOKEN
	}

	// Check if model is set, if not set to default
	if userParams.Request.Model == "" {
		params.Request.Model = "gpt-3.5-turbo"
	} else {
		params.Request.Model = userParams.Request.Model
	}

	// Check if KeepMessageHistory is set, if not set to default
	if userParams.KeepMessageHistory {
		params.KeepMessageHistory = true
	}

	// Check if temperature is set, if not set to default
	if userParams.StripNewline {
		params.StripNewline = true
	}

	return params, nil
}

func (params *Params) ClearHistory(msg string) {
	// Go over all params.MessageHistory and remove all messages
	for i := 0; i < len(params.MessageHistory); i++ {
		params.MessageHistory = append(params.MessageHistory[:i], params.MessageHistory[i+1:]...)
	}
}
