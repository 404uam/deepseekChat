package ai

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"

	"github.com/openai/openai-go"
)

func DoAi(client openai.Client, ch chan string, prompt string, previous string, chatModel ChatModel) string {

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	defer close(ch)

	chatCompletion, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(previous + "\n" + prompt),
			},
			Seed:  openai.Int(1),
			Model: chatModel,
		})

	if err != nil {
		panic(err.Error())
	}

	response := chatCompletion.Choices[0].Message.Content
	println(response)
	ch <- response
	return response
}

func DoAiWithStreaming(client openai.Client, ch chan string, pressed time.Time, markdown *widget.RichText, label *container.Scroll, prompt string, previous string, chatModel ChatModel) string {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	chatCompletionStream := client.Chat.Completions.NewStreaming(ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(previous + "\n" + prompt),
			},
			Temperature: openai.Float(1.3),
			Seed:        openai.Int(1),
			Model:       chatModel,
		})
	acc := openai.ChatCompletionAccumulator{}
	timeElapsed := time.Duration(0)
	label.Show()
	markdown.Hide()
	defer close(ch)
	for chatCompletionStream.Next() {
		currChunk := chatCompletionStream.Current()
		if timeElapsed == 0 {
			timeElapsed = time.Since(pressed)
		}
		acc.AddChunk(currChunk)

		// When this fires, the current chunk value will not contain content data
		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
			println()
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
			println()
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
			println()
		}

		// It's best to use chunks after handling JustFinished events
		if len(currChunk.Choices) > 0 {
			content := currChunk.Choices[0].Delta.JSON.Content
			if currChunk.Choices[0].FinishReason != "stop" {
				var tempWords string
				err := json.Unmarshal([]byte(content.Raw()), &tempWords)
				if err != nil {
					println("got err:", err.Error())
					println(content.Raw())
					tempWords = ""
				}
				println(tempWords)
				ch <- tempWords
			} else {
				ch <- ""
			}
		}

	}
	if err := chatCompletionStream.Err(); err != nil {
		markdown.AppendMarkdown(err.Error())
		label.Hide()
		label.Refresh()
		markdown.Show()
		markdown.Refresh()
		return ""
	}

	println("Total Tokens:", acc.Usage.TotalTokens)
	println("Finish Reason:", acc.Choices[0].FinishReason)
	println("Time from button press to first response: ", timeElapsed.String())
	markdown.AppendMarkdown(acc.Choices[0].Message.Content + "<br><br>")
	markdown.AppendMarkdown("<br><br>")
	label.Hide()
	label.Refresh()
	markdown.Show()
	markdown.Refresh()
	return acc.Choices[0].Message.Content
}

func DoAssistantAiWithStreaming(client openai.Client, ch chan string, pressed time.Time, markdown *widget.RichText, label *container.Scroll, prompt string, previous string, chatModel ChatModel, systemPrompt string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	chatCompletionStream := client.Chat.Completions.NewStreaming(ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.DeveloperMessage(systemPrompt),
				openai.UserMessage(previous + "\n" + prompt),
			},
			Temperature: openai.Float(1.3),
			Seed:        openai.Int(1),
			Model:       chatModel,
		})
	acc := openai.ChatCompletionAccumulator{}
	timeElapsed := time.Duration(0)
	label.Show()
	markdown.Hide()
	defer close(ch)
	for chatCompletionStream.Next() {
		currChunk := chatCompletionStream.Current()
		if timeElapsed == 0 {
			timeElapsed = time.Since(pressed)
		}
		acc.AddChunk(currChunk)

		// When this fires, the current chunk value will not contain content data
		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
			println()
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
			println()
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
			println()
		}

		// It's best to use chunks after handling JustFinished events
		if len(currChunk.Choices) > 0 {
			content := currChunk.Choices[0].Delta.JSON.Content
			if currChunk.Choices[0].FinishReason != "stop" {
				var tempWords string
				err := json.Unmarshal([]byte(content.Raw()), &tempWords)
				if err != nil {
					println("got err:", err.Error())
					println(content.Raw())
					tempWords = ""
				}
				println(tempWords)
				ch <- tempWords
			} else {
				ch <- ""
			}
		}

	}
	if err := chatCompletionStream.Err(); err != nil {
		markdown.AppendMarkdown(err.Error())
		label.Hide()
		label.Refresh()
		markdown.Show()
		markdown.Refresh()
		return ""
	}

	println("Total Tokens:", acc.Usage.TotalTokens)
	println("Finish Reason:", acc.Choices[0].FinishReason)
	println("Time from button press to first response: ", timeElapsed.String())
	markdown.AppendMarkdown(acc.Choices[0].Message.Content + "<br><br>")
	markdown.AppendMarkdown("<br><br>")
	label.Hide()
	label.Refresh()
	markdown.Show()
	markdown.Refresh()
	return acc.Choices[0].Message.Content
}
