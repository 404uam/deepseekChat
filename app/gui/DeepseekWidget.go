package gui

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)
import "deepseekChat/m/app/ai"

func NewDeepseekChatWidget(apiKey string) *AiWidget {
	chatModel := ai.DeepseekChat
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com"),
	)
	return NewAiWidget(client, chatModel)
}

func NewDeepseekReasonerWidget(apiKey string) *AiWidget {
	chatModel := ai.DeepseekChat

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com"),
	)
	return NewAiWidget(client, chatModel)
}
