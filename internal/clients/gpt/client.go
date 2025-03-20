package gpt

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client openai.Client
}

func NewClient(apiKey string, baseURL string) *Client {
	defaultConfig := openai.DefaultConfig(apiKey)
	defaultConfig.BaseURL = baseURL

	gptClient := openai.NewClientWithConfig(defaultConfig)

	return &Client{
		client: *gptClient,
	}
}

func (c *Client) GenerateDiet(prompt string) (string, error) {
	response, err := c.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
