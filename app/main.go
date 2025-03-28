package main

import (
	"deepseekChat/m/app/gui"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"os"
)

type Config struct {
	DeepseekAPIKey string `json:"API_KEY"`
	OpenAiAPIKey   string `json:"OPENAI_API_KEY"`
}

func main() {
	file, err := os.ReadFile("app/config/config.json")
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Error parsing config file:", err)
	}
	if config.DeepseekAPIKey == "" {
		log.Fatal("API Key not found in config")
	}

	myapp := app.New()
	window := myapp.NewWindow("CHatter")
	myCanvas := window.Canvas()

	textWidget := widget.NewRichTextFromMarkdown("aiResponse")
	textWidget.Wrapping = fyne.TextWrapWord
	textWidget.Resize(fyne.NewSize(50, 50))

	markdown := widget.NewRichTextFromMarkdown("")
	markdown.Wrapping = fyne.TextWrapWord

	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord

	labelScrollContainer := container.NewScroll(label)
	labelScrollContainer.Hide()
	activity := widget.NewActivity()
	activity.Hide()

	boundString := binding.NewString()
	label.Bind(boundString)
	/*	submitButton := widget.NewButton("ask away!", func() {
			ch := make(chan string, 1)
			inputText := input.Text
			log.Println(inputText)
			input.Disable()
			input.SetText("")
			activity.Start()
			activity.Show()
			prevString, _ := boundString.Get()
			go ai.DoAi(client, ch, inputText, prevString, openai.ChatModelGPT4oMini)
			aiResponse := <-ch
			builder2 := strings.Builder{}
			builder2.WriteString(prevString + "\n\n" + aiResponse)
			markdown.AppendMarkdown(builder2.String() + "\n\n")
			boundString.Set(builder2.String())
			markdownScrollContainer.ScrollToBottom()
			activity.Stop()
			activity.Hide()
			input.Enable()
			inputAndButtons.Refresh()
		})
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
				openai.ChatModelGPT3_5Turbo,
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
				inputAndButtons.Refresh()
			}()
		})*/

	whiteRectangle := canvas.NewRectangle(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 75,
	})
	whiteRectangle.Resize(fyne.NewSize(125, 200))
	whiteRectangle.SetMinSize(whiteRectangle.Size())

	tabs := container.NewAppTabs(
		container.NewTabItem("ChatGPT", widget.NewLabel("LL")),
		container.NewTabItem("Deepseekchat", gui.NewDeepseekWidget(config.DeepseekAPIKey)),
	)

	myCanvas.SetContent(tabs)

	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
