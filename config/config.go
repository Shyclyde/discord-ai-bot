package config

import (
	"encoding/json"
	"log"
	"os"
)

type AppConfig struct {
	Bot struct {
		Name         string `json:"name"`
		BotAdminRole string `json:"botAdminRole"`
	} `json:"bot"`

	Discord struct {
		OtherBotID    string `json:"otherBotID"`
		TextChannelID string `json:"textChannelID"`
	} `json:"discord"`

	OpenAI struct {
		Text struct {
			Enabled          bool    `json:"enabled"`
			API_URL          string  `json:"apiURL"`
			Model            string  `json:"model"`
			MaxTokens        int     `json:"maxTokens"`
			TopP             float32 `json:"topP"`
			Temperature      float32 `json:"temperature"`
			Completions      int     `json:"completions"`
			FrequencyPenalty float32 `json:"frequencyPenalty"`
			PresencePenalty  float32 `json:"presencePenalty"`
			StopValue        string  `json:"stopValue"`
			ContextLength    int     `json:"contextLength"`
		} `json:"text"`
		Image struct {
			Enabled     bool   `json:"enabled"`
			API_URL     string `json:"apiURL"`
			Completions int    `json:"completions"`
			ImageSize   string `json:"imageSize"`
		} `json:"image"`
	} `json:"openAI"`

	GameServer struct {
		ServerName string `json:"serverName"`
		Games      []struct {
			Name           string `json:"name"`
			ProcessName    string `json:"processName"`
			MemoryEstimate int    `json:"memoryEstimate"`
		}
	} `json:"gameServer"`
}

var (
	Config AppConfig
)

func init() {
	loadConfig(&Config)
}

func loadConfig(config *AppConfig) {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Can't open config.json file:", err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatalln("Can't decode config.json file:", err)
	}
}
