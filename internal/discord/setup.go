package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func LoadServerConfig(s *discordgo.Session) {
	initIntents(s)
	initHandlers(s)
	loadDiscordAppCommands(s)
}

func initIntents(s *discordgo.Session) {
	log.Printf("Loading Discord intents...\n")

	s.Identify.Intents = discordgo.IntentGuildMessages |
		discordgo.IntentGuilds |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildEmojis |
		discordgo.IntentGuildIntegrations |
		discordgo.IntentMessageContent
}

func initHandlers(s *discordgo.Session) {
	log.Printf("Binding Discord event handlers...\n")
	s.AddHandler(MessageHandler)
	//s.AddHandler(eventhandlers.ReactionCreateHandler) // This seems to be deprecated now?
	s.AddHandler(EmojiUpdateHandler)
	s.AddHandler(InteractionHandler)
}
