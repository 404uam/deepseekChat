package gui

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewAssistantWidget(apiKey string) *AiWidget {
	chatModel := openai.ChatModelGPT4oMini
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return NewAssistantAiWidget(client, chatModel)
}
