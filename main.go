package main

import (
	"context"

	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithAPIKey(""), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL("https://api.deepseek.com"),
	)

	chatCompletion, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello say hi"),
			}),
			Model: openai.F(openai.ChatModelDeepSeek),
		})

	/* chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	}) */
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)

}
