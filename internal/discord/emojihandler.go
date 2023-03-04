package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func EmojiUpdateHandler(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {
	/*
		there seems to be some limitations with this event.
		we only get a list of guild emojis back, no way to tell if an emoji
		was added or removed without keeping some previous state somewhere.
		meh.
	*/

	latestemoji := e.Emojis[len(e.Emojis)-1]
	fmt.Printf("Emojis updated, current latest emoji: %s\n", latestemoji.Name)
}
