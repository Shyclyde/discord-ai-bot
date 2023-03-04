package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/shyclyde/discord-ai-bot/config"
)

func CheckIsBotAdmin(session *discordgo.Session, member *discordgo.Member, gid string) bool {
	for _, roleid := range member.Roles {
		role, err := session.State.Role(gid, roleid)
		if err != nil {
			log.Printf("Error, couldn't retrieve role %s, %s\n", roleid, err)
			return false
		}
		if role.Name == config.Config.Bot.BotAdminRole {
			return true
		}
	}
	return false
}
