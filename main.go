package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shyclyde/discord-ai-bot/config"
	"github.com/shyclyde/discord-ai-bot/internal/discord"
	"github.com/shyclyde/discord-ai-bot/internal/stats"
)

var (
	token string
)

func init() {
	token = os.Getenv("DISCORD_BOT_TOKEN")

	err := checkEnv()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func checkEnv() error {
	if token == "" {
		return errors.New("no DISCORD_BOT_TOKEN environment variable found")
	}
	return nil
}

func main() {
	log.Printf("Creating Discord bot session...\n")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Error: couldn't create Discord session, %v\n", err)
		os.Exit(1)
	}

	log.Printf("Opening Discord websocket connection...\n")
	err = bot.Open()
	if err != nil {
		log.Printf("Error: couldn't open Discord websocket connection, %v\n", err)
		os.Exit(1)
	}
	defer bot.Close()

	discord.LoadServerConfig(bot)
	go stats.StartUpdateTick(bot)

	// Bot is ready, keep online till the bot is killed
	log.Printf("%s is ready to serve. (ctrl+c to exit)\n", config.Config.Bot.Name)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Printf("Shutting down...\n")

	discord.RemoveDiscordAppCommands(bot)
	stats.RemoveStats(bot)

	log.Printf("%s is flying away.\n", config.Config.Bot.Name)
}
