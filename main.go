package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	openai "github.com/sashabaranov/go-openai"
)

var apiKey = ""

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	// systray.SetIcon(getIcon("default")) // Set your default icon here
	systray.SetTitle("AI Helper")
	systray.SetTooltip("AI Helper")

	mDefault := systray.AddMenuItem("Proofread Clipboard", "Send clipboard content to OpenAI for proofreading")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mDefault.ClickedCh:
				handleDefaultAction()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Clean up here
}

func handleDefaultAction() {
	content, err := clipboard.ReadAll()
	if err != nil {
		log.Println("Failed to read clipboard content:", err)
		return
	}

	// systray.SetIcon(getIcon("loading")) // Set your loading icon here

	response, err := sendToOpenAI(content)
	if err != nil {
		log.Println("Failed to send to OpenAI:", err)
		return
	}

	err = clipboard.WriteAll(response)
	if err != nil {
		log.Println("Failed to write to clipboard:", err)
		return
	}

	// systray.SetIcon(ge2tIcon("default")) // Revert to default icon
}
func sendToOpenAI(content string) (string, error) {
	client := openai.NewClient(apiKey)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "Please fix any typos and grammar in the language that is provided in. Don't rewrite slangs, only respond with the corrected text:\n" + content,
			},
		},
		MaxTokens: 100,
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response choices available")
}
func getIcon(status string) []byte {
	// Load your icons here. You can return different icons based on the status.
	// For example:
	// - "default": Default state icon
	// - "loading": Loading state icon
	// You can embed your icons as byte arrays or load them from files.

	iconPath := "./default/icon.ico"
	if status == "loading" {
		iconPath = "./loading/icon.ico"
	}

	icon, err := ioutil.ReadFile(iconPath)
	if err != nil {
		log.Fatal(err)
	}
	return icon
}
