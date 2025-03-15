package main

import (
	"context"
	"time"

	"github.com/openai/openai-go"
)

func DoAi(client *openai.Client, ch chan string, prompt string, previous string) string {

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	chatCompletion, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(previous + "\n" + prompt),
			}),
			Seed:  openai.Int(1),
			Model: openai.F(DeepseekChat),
		})

	if err != nil {
		panic(err.Error())
	}

	response := chatCompletion.Choices[0].Message.Content
	println(response)
	ch <- response
	return response
}
