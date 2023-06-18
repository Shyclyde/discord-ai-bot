package discord

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/shyclyde/discord-ai-bot/config"
	"github.com/shyclyde/discord-ai-bot/internal/gameservers"
	"github.com/shyclyde/discord-ai-bot/internal/openai"
	"github.com/shyclyde/discord-ai-bot/pkg/utils"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.ChannelID != config.Config.Discord.TextChannelID {
		return
	} else if strings.Contains(strings.ToLower(m.Content), strings.ToLower(config.Config.Bot.Name)) {
		msg := strings.ToLower(m.Content)

		otherbot, err := s.GuildMember(m.GuildID, config.Config.Discord.OtherBotID)
		if err != nil {
			log.Printf("Can't get other bot, %v\n", err)
			return
		}
		otherbotname := strings.ToLower(otherbot.User.Username)

		if msg == "am i a bot admin?" {
			msgCheckSelfAdmin(s, m)
		} else if strings.Contains(msg, fmt.Sprintf("restart %s", otherbotname)) {
			msgRestartOtherBot(s, m)
		} else if strings.Contains(msg, "game servers") {
			msgCheckGameServers(s, m)
		} else if strings.Contains(msg, "how's the server memory?") {
			msgCheckServerMemory(s, m)
		} else if strings.Contains(msg, "serve me ") {
			log.Println("Handling image request")
			msgOpenAIImage(s, m)
		} else if config.Config.OpenAI.Text.Enabled {
			msgOpenAIText(s, m)
		}
	}
}

func msgOpenAIText(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("OpenAI text request: %s\n", m.Content)
	reply, err, max_reached := openai.GenerateText(m.Content, m.Author.Username)
	if err != nil {
		log.Printf("Error trying to generate OpenAI text: %s\n", err)
		s.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	log.Printf("OpenAI text completion: %s\n", reply)
	s.ChannelMessageSend(m.ChannelID, reply)

	if max_reached {
		s.ChannelMessageSend(m.ChannelID, "Oops, I ran out of breath...")
	}
}

func msgOpenAIImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	prompt := strings.SplitAfter(strings.ToLower(m.Content), "serve me ")

	if len(prompt) < 2 {
		log.Printf("Invalid image request: %s", m.Content)
		s.ChannelMessageSend(m.ChannelID, "There's nothing to serve.")
		return
	}

	req := prompt[1]
	if req == "" || req == " " {
		log.Printf("Blank image request")
		s.ChannelMessageSend(m.ChannelID, "There's nothing to serve.")
		return
	}

	log.Printf("OpenAI text request: %s\n", req)
	imageURL, err := openai.GenerateImage(req)
	if err != nil {
		log.Printf("Error trying to generate OpenAI image: %s", err)
		s.ChannelMessageSend(m.ChannelID, "I'm sorry, that's unfortunately not something I can serve.")
		return
	}

	log.Printf("OpenAI image URL: %s\n", imageURL)
	filename, err := openai.DownloadOpenAIImage(imageURL)
	if err != nil {
		log.Printf("Error trying to download OpenAI image: %s", err)
		s.ChannelMessageSend(m.ChannelID, "I'm sorry, I couldn't download that.")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error trying to open OpenAI image: %s\n", err)
		s.ChannelMessageSend(m.ChannelID, "I'm sorry, I couldn't open the file.")
		return
	}

	_, err = s.ChannelFileSend(m.ChannelID, req+".png", file)
	if err != nil {
		log.Printf("Error trying to send OpenAI image to Discord: %s\n", err)
		s.ChannelMessageSend(m.ChannelID, "I'm sorry, I couldn't upload the image.")
		return
	}
}

func msgCheckServerMemory(s *discordgo.Session, m *discordgo.MessageCreate) {
	free := fmt.Sprintf("%d MB\n", utils.GetSysFreeMemory()/1024/1024)
	total := fmt.Sprintf("%d MB\n", utils.GetSysTotalMemory()/1024/1024)
	s.ChannelMessageSend(m.ChannelID, free)
	s.ChannelMessageSend(m.ChannelID, total)
}

func msgCheckGameServers(s *discordgo.Session, m *discordgo.MessageCreate) {
	statusMessage := "Game Server Status:\n"
	for _, game := range gameservers.GetGameServerStatus() {

		statusMessage += fmt.Sprintf("%s %s\n", game.Active, game.Name)
	}
	s.ChannelMessageSend(m.ChannelID, statusMessage)
}

func msgRestartOtherBot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !CheckIsBotAdmin(s, m.Member, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, "I bow to no peasant.")
		return
	}

	if !utils.HandleProcess("spoopy", "restart") {
		s.ChannelMessageSend(m.ChannelID, "Something went wroing trying to restart Spoopy...")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Spoopy has new life.")
	}
}

func msgCheckSelfAdmin(s *discordgo.Session, m *discordgo.MessageCreate) {
	if CheckIsBotAdmin(s, m.Member, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, "Yes, at your service.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "No, go away.")
	}
}
