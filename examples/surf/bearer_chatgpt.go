package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/x0xO/hg/surf"
)

func main() {
	reply, err := Completions("who is mr. root?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Join(strings.Fields(reply), " "))
}

const (
	BASEURL = "https://api.openai.com/v1/"
	APIKEY  = ""
)

type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        uint    `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	BestOf           int     `json:"best_of"`
}

func Completions(msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        1000,
		Temperature:      0.5,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		BestOf:           5,
	}

	opt := surf.NewOptions().BearerAuth(APIKEY)

	r, err := surf.NewClient().SetOptions(opt).Post(BASEURL+"completions", requestBody).Do()
	if err != nil {
		return "", err
	}

	var gptResponseBody ChatGPTResponseBody

	err = r.Body.JSON(&gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Text
	}

	return reply, nil
}
