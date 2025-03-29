package gui

import (
	"deepseekChat/m/app/ai"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"log"
	"strings"
	"time"
)

type ChatGPTWidget struct {
	widget.BaseWidget
	ChatModel            openai.ChatModel
	textInputWidget      *fyne.Container
	markdownOrLabelStack *fyne.Container
}

func NewChatGPTWidget(apiKey string) *ChatGPTWidget {
	chatModel := openai.ChatModelGPT4oMini
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	activity := widget.NewActivity()
	activity.Hide()
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your prompt...")
	inputActivity := container.NewBorder(nil, nil, activity, nil, input)

	markdown := widget.NewRichTextFromMarkdown("")
	markdown.Wrapping = fyne.TextWrapWord

	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord

	markdownScrollContainer := container.NewScroll(markdown)
	labelScrollContainer := container.NewScroll(label)

	boundString := binding.NewString()
	label.Bind(boundString)

	submitWithStreamingButton := widget.NewButton("stream away!", func() {
		buttonPressed := time.Now()
		inputText := input.Text
		ch2 := make(chan string, 1)
		log.Println(inputText)
		input.Disable()
		input.SetText("")
		activity.Start()
		activity.Show()
		prevString, _ := boundString.Get()
		go ai.DoAiWithStreaming(
			client,
			ch2,
			buttonPressed,
			markdown,
			labelScrollContainer,
			inputText,
			prevString,
			chatModel,
		)

		go func() {
			builder2 := strings.Builder{}
			builder2.WriteString(prevString)
			for aiResponse := range ch2 {
				builder2.WriteString(aiResponse)
				boundString.Set(builder2.String())
				labelScrollContainer.ScrollToBottom()
				markdownScrollContainer.ScrollToBottom()
			}
			builder2.WriteString("\n\n")
			boundString.Set(builder2.String())
			markdownScrollContainer.ScrollToBottom()
			activity.Stop()
			activity.Hide()
			input.Enable()
			inputActivity.Refresh()
		}()
	})

	chatGPTWidget := &ChatGPTWidget{
		ChatModel:            chatModel,
		textInputWidget:      container.NewVBox(inputActivity, submitWithStreamingButton),
		markdownOrLabelStack: container.NewStack(markdownScrollContainer, labelScrollContainer),
	}
	chatGPTWidget.ExtendBaseWidget(chatGPTWidget)
	return chatGPTWidget
}

func (d *ChatGPTWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(nil, d.textInputWidget, nil, nil, d.markdownOrLabelStack),
	)
}
