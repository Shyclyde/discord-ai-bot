package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/shyclyde/discord-ai-bot/config"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: fmt.Sprintf("Say hi to %s!", config.Config.Bot.Name),
		},
		{
			Name:        "gameserver",
			Description: "Adminstration for game servers",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "game",
					Description: "Which game server to administrate",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "process-action",
					Description: "Process action option (restart, start, stop)",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "restart",
							Value: "restart",
						},
						{
							Name:  "start",
							Value: "start",
						},
						{
							Name:  "stop",
							Value: "stop",
						},
					},
					Required: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello":      HelloHandler,
		"gameserver": GameServerHandler,
	}
)

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handlerFunc, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handlerFunc(s, i)
	}
}

func RemoveDiscordAppCommands(s *discordgo.Session) error {
	log.Printf("Cleaning up app commands...\n")

	appCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Printf("Can't retrieve app commands to delete: %s\n", err)
	}

	for _, cmd := range appCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
	return nil
}

func loadDiscordAppCommands(s *discordgo.Session) error {
	log.Printf("Registering slash commands...\n")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	return nil
}
