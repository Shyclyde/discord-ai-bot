package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/shyclyde/discord-ai-bot/config"
	"github.com/shyclyde/discord-ai-bot/pkg/utils"
)

type OpenAITextResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type OpenAITextRequest struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Top_P            float32 `json:"top_p"`
	Temperature      float32 `json:"temperature"`
	Completions      int     `json:"n"`
	StopValue        string  `json:"stop"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}

type OpenAIImageRequest struct {
	Prompt      string `json:"prompt"`
	Completions int    `json:"n"`
	Size        string `json:"size"`
}

type OpenAIImageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

var (
	openai_api_key     string
	gptContextMessages []string
)

func init() {
	openai_api_key = os.Getenv("OPENAI_API_KEY")
	err := checkEnv()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func checkEnv() error {
	if openai_api_key == "" {
		return errors.New("no OPENAI_API_KEY environment variable found")
	}
	return nil
}

func callOpenAIAPI(body []byte, API_URL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openai_api_key))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func GenerateText(prompt string, requestor string) (string, error, bool) {
	max_reached := false
	training, err := os.ReadFile("openai_training.txt")
	if err != nil {
		log.Fatal("cannot read training file", err)
	}

	GPTContextString := ""
	for _, item := range gptContextMessages {
		GPTContextString += item
	}

	full_prompt := fmt.Sprintf(
		"%s\n%s%s: %s\n%s:",
		string(training),
		GPTContextString,
		requestor,
		prompt,
		config.Config.Bot.Name,
	)

	payload := OpenAITextRequest{
		Model:            config.Config.OpenAI.Text.Model,
		Prompt:           full_prompt,
		MaxTokens:        config.Config.OpenAI.Text.MaxTokens,
		Top_P:            config.Config.OpenAI.Text.TopP,
		Temperature:      config.Config.OpenAI.Text.Temperature,
		Completions:      config.Config.OpenAI.Text.Completions,
		PresencePenalty:  config.Config.OpenAI.Text.PresencePenalty,
		FrequencyPenalty: config.Config.OpenAI.Text.FrequencyPenalty,
		StopValue:        config.Config.OpenAI.Text.StopValue,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err, max_reached
	}

	respBody, err := callOpenAIAPI(body, config.Config.OpenAI.Text.API_URL)
	if err != nil {
		return "", err, max_reached
	}

	var respObj OpenAITextResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return "", err, max_reached
	}

	if len(respObj.Choices) < 1 {
		return "", fmt.Errorf("no responses from OpenAI text completion"), max_reached
	}

	if respObj.Choices[0].FinishReason == "length" {
		max_reached = true
		log.Printf("Warning! OpenAI max_length reached in response\n")
	}

	gptContextMessages = append(gptContextMessages, fmt.Sprintf("%s: %s\n%s: %s\n",
		requestor,
		prompt,
		config.Config.Bot.Name,
		respObj.Choices[0].Text,
	))
	if len(gptContextMessages) > config.Config.OpenAI.Text.ContextLength {
		gptContextMessages = gptContextMessages[1 : config.Config.OpenAI.Text.ContextLength+1]
	}

	utils.LogMemUsage()
	return respObj.Choices[0].Text, nil, max_reached
}

func GenerateImage(prompt string) (string, error) {
	payload := OpenAIImageRequest{
		Prompt:      prompt,
		Completions: config.Config.OpenAI.Image.Completions,
		Size:        config.Config.OpenAI.Image.ImageSize,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	respBody, err := callOpenAIAPI(body, config.Config.OpenAI.Image.API_URL)
	if err != nil {
		return "", err
	}

	var respObj OpenAIImageResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return "", err
	}

	if len(respObj.Data) < 1 {
		return "", fmt.Errorf("no responses from OpenAI image completion")
	}

	return respObj.Data[0].URL, nil
}
