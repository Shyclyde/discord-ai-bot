package gameservers

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type gameServer struct {
	Name           string `json:"name"`
	ProcessName    string `json:"processName"`
	MemoryEstimate int    `json:"memoryEstimate"`
	IsActive       bool
}

type serverConfig struct {
	ServerName  string       `json:"serverName"`
	GameServers []gameServer `json:"gameServers"`
}

type GameStatus struct {
	Name   string
	Active string
}

var server serverConfig

func init() {
	//LoadAllGameServers()
}

func LoadAllGameServers() {
	log.Printf("Loading server_config.json file...\n")
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error when opening config.json:", err)
	}

	err = json.Unmarshal(data, &server)
	if err != nil {
		log.Fatal("Error unmarhsaling config.json:", err)
	}
}

func CheckAllGameServers() {
	// Work in progress
}

func GetGameServerStatus() []GameStatus {
	var games []GameStatus
	for _, game := range server.GameServers {
		active := ":x:"
		if game.IsActive {
			active = ":white_check_mark:"
		}
		games = append(games, GameStatus{Name: game.Name, Active: active})
	}
	return games
}

func CheckGameServer(gameCheck string) bool {
	for _, game := range server.GameServers {
		if strings.EqualFold(gameCheck, game.Name) {
			return game.IsActive
		}
	}
	return false
}
