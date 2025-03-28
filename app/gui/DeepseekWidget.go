package gui

import (
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
import "deepseekChat/m/app/ai"

type DeepseekWidget struct {
	widget.BaseWidget
	ChatModel            ai.ChatModel
	textInputWidget      *fyne.Container
	markdownOrLabelStack *fyne.Container
}

func NewDeepseekWidget(apiKey string) *DeepseekWidget {
	chatModel := ai.DeepseekChat
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com"),
	)
	activity := widget.NewActivity()
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

	deepseekWidget := &DeepseekWidget{
		ChatModel:            chatModel,
		textInputWidget:      container.NewVBox(inputActivity, submitWithStreamingButton),
		markdownOrLabelStack: container.NewStack(markdownScrollContainer, labelScrollContainer),
	}
	deepseekWidget.ExtendBaseWidget(deepseekWidget)
	return deepseekWidget
}

func (d DeepseekWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(nil, d.textInputWidget, nil, nil, d.markdownOrLabelStack),
	)
}
