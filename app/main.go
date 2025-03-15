package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"image/color"
	"log"
	"os"
	"strings"
)

type Config struct {
	APIKey       string `json:"API_KEY"`
	OpenAiAPIKey string `json:"OPENAI_API_KEY"`
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
	if config.APIKey == "" {
		log.Fatal("API Key not found in config")
	}

	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL("https://api.deepseek.com"),
	)

	myapp := app.New()
	window := myapp.NewWindow("CHatter")
	myCanvas := window.Canvas()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your prompt...")

	textBox := canvas.NewText("aiResponse", color.Black)

	textBox.Resize(fyne.NewSize(200, 200))
	textWidget := widget.NewRichTextFromMarkdown("aiResponse")
	textWidget.Wrapping = fyne.TextWrapWord
	textWidget.Resize(fyne.NewSize(50, 50))

	markdown := widget.NewRichTextFromMarkdown("")
	markdown.Wrapping = fyne.TextWrapWord
	scrollContainer := container.NewScroll(markdown)

	ch := make(chan string, 1)
	boundString := binding.NewString()
	submitButton := widget.NewButton("ask away!", func() {
		inputText := input.Text
		log.Println(inputText)
		input.Disable()
		input.SetText("")
		prevString, _ := boundString.Get()
		go DoAi(client, ch, inputText, prevString)
		aiResponse := <-ch
		builder2 := strings.Builder{}
		builder2.WriteString(prevString + "\n\n" + aiResponse)
		markdown.AppendMarkdown(builder2.String() + "\n\n")
		boundString.Set(builder2.String())
		scrollContainer.ScrollToBottom()
		input.Enable()
	})

	textInputWidget := container.NewVBox(input, submitButton)

	whiteRectangle := canvas.NewRectangle(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 75,
	})
	whiteRectangle.Resize(fyne.NewSize(125, 200))
	whiteRectangle.SetMinSize(whiteRectangle.Size())

	content := container.New(layout.NewHBoxLayout(), textWidget, whiteRectangle)
	content = container.NewBorder(nil, textInputWidget, whiteRectangle, nil, scrollContainer)

	myCanvas.SetContent(content)

	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
