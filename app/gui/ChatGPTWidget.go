package gui

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewChatGPTWidget(apiKey string) *AiWidget {
	chatModel := openai.ChatModelGPT4oMini
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return NewAiWidget(client, chatModel)
}
