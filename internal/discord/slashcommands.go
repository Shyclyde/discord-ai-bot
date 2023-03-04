package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shyclyde/discord-ai-bot/pkg/utils"
)

func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "hello!",
		},
	})
}

func GameServerHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !CheckIsBotAdmin(s, i.Member, i.GuildID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I bow to no peasant. Admins only.",
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options
	game := options[0].StringValue()
	action := options[1].StringValue()

	if game != "minecraft" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Only supporting Minecraft during testing",
			},
		})
		return
	}

	msgformat := fmt.Sprintf("Trying to %s %s...\n", action, game)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msgformat,
		},
	})

	if !utils.HandleProcess(game, "restart") {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Something went wroing trying to restart %s.\n", game),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%s has been restarted.\n", game),
			},
		})
	}
}
