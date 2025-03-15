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

	label := widget.NewLabel("AI Response:")

	ch := make(chan string, 1)
	boundString := binding.NewString()
	label.Bind(boundString)
	submitButton := widget.NewButton("ask away!", func() {
		log.Println(input.Text)

		prevString, _ := boundString.Get()
		go DoAi(client, ch, input.Text, prevString)
		aiResponse := <-ch
		builder2 := strings.Builder{}
		builder2.WriteString(prevString + "\n\n" + aiResponse)
		boundString.Set(builder2.String())
	})

	textInputWidget := container.NewVBox(input, submitButton)

	label.Wrapping = fyne.TextWrapBreak
	label.Resize(fyne.NewSize(200, 200))
	whiteRectangle := canvas.NewRectangle(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 100,
	})
	whiteRectangle.Resize(fyne.NewSize(125, 200))
	whiteRectangle.SetMinSize(whiteRectangle.Size())

	scrollContainer := container.NewScroll(label)

	content := container.New(layout.NewHBoxLayout(), textWidget, whiteRectangle)
	content = container.NewBorder(nil, textInputWidget, whiteRectangle, nil, scrollContainer)

	myCanvas.SetContent(content)

	window.Resize(fyne.NewSize(800, 600))
	/*
		go func() {
			time.Sleep(time.Second)
			previousPrompt, _ := boundString.Get()
			DoAi(client, ch, "what about the second?", previousPrompt)
			blah := <-ch
			println("got here foobar")

			builder := strings.Builder{}
			builder.WriteString(previousPrompt + "\n\n" + blah)

			boundString.Set(builder.String())

			previousPrompt2, _ := boundString.Get()
			DoAi(client, ch, "what about mount lougheed?", previousPrompt2)
			blah2 := <-ch
			println("got here foobar")

			builder2 := strings.Builder{}
			builder2.WriteString(previousPrompt + "\n\n" + blah2)

			boundString.Set(builder2.String())
		}()*/

	window.ShowAndRun()
}

func buildUI() {
	myapp := app.New()
	window := myapp.NewWindow("CHatter")
	myCanvas := window.Canvas()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your prompt...")
	textInputWidget := container.NewVBox(input, widget.NewButton("print", func() {
		log.Println(input.Text)
	}))

	textBox := canvas.NewText("aiResponse", color.Black)

	textBox.Resize(fyne.NewSize(200, 200))
	textWidget := widget.NewRichTextFromMarkdown("aiResponse")
	textWidget.Wrapping = fyne.TextWrapWord
	textWidget.Resize(fyne.NewSize(50, 50))

	whiteRectangle := canvas.NewRectangle(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 100,
	})
	whiteRectangle.Resize(fyne.NewSize(200, 200))
	whiteRectangle.SetMinSize(whiteRectangle.Size())
	textBox.Move(fyne.NewPos(0, 0))
	content := container.New(layout.NewHBoxLayout(), textWidget, whiteRectangle)
	content = container.NewBorder(nil, textInputWidget, whiteRectangle, nil, textWidget)

	myCanvas.SetContent(content)

	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
