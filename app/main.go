package main

import (
	"deepseekChat/m/app/gui"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	tabs := container.NewAppTabs(
		container.NewTabItem("ChatGPT", gui.NewChatGPTWidget(config.OpenAiAPIKey)),
		container.NewTabItem("Deepseekchat", gui.NewDeepseekWidget(config.DeepseekAPIKey)),
	)

	myCanvas.SetContent(tabs)

	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
