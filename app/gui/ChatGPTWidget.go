package gui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
)

type ChatGPTWidget struct {
	widget.BaseWidget
	ChatModel openai.ChatModel
}
