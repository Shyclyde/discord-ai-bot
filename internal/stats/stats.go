package stats

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	UpdateTickSeconds = 10
)

func StartUpdateTick(s *discordgo.Session) {
	for {
		for _, guild := range s.State.Guilds {
			UpdateGuild(s, guild)
		}
		time.Sleep(UpdateTickSeconds * time.Second)
	}
}

func RemoveStats(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		log.Printf("Removing stats from %s...\n", guild.Name)
		RemoveStatsFromGuild(s, guild)
	}
}

func RemoveStatsFromGuildCategory(s *discordgo.Session, g *discordgo.Guild, parent_id string) error {
	for _, ch := range g.Channels {
		if ch.ParentID == parent_id {
			_, err := s.ChannelDelete(ch.ID)
			if err != nil {
				log.Printf("Couldn't delete channel %s: %v\n", ch.Name, err)
			}
		}
	}

	ch, err := s.ChannelDelete(parent_id)
	if err != nil {
		log.Printf("Couldn't delete stats category channel %s: %v\n", ch.Name, err)
	}

	return nil
}

func GetStatsCategoryCount(s *discordgo.Session, g *discordgo.Guild) int {
	count := 0
	for _, ch := range g.Channels {
		if ch.Name == fmt.Sprintf("%s Stats", g.Name) {
			count++
		}
	}
	return count
}

func GetStatsCategoryID(s *discordgo.Session, g *discordgo.Guild) (string, error) {
	count := GetStatsCategoryCount(s, g)
	if count > 1 {
		RemoveStatsFromGuild(s, g)
		CreateStatsOnGuild(s, g)
	} else if count == 0 {
		CreateStatsOnGuild(s, g)
	}

	ch_id := ""
	for _, ch := range g.Channels {
		if ch.Name == fmt.Sprintf("%s Stats", g.Name) {
			ch_id = ch.ID
		}
	}

	if ch_id == "" {
		err_msg := fmt.Sprintf("failure getting of stats for %s\n", g.Name)
		return "", errors.New(err_msg)
	}

	return ch_id, nil
}

func RemoveStatsFromGuild(s *discordgo.Session, g *discordgo.Guild) {
	for _, ch := range g.Channels {
		if ch.Name == fmt.Sprintf("%s Stats", g.Name) {
			RemoveStatsFromGuildCategory(s, g, ch.ID)
		}
	}
}

func CreateStatsOnGuild(s *discordgo.Session, g *discordgo.Guild) {
	cat_data := discordgo.GuildChannelCreateData{
		Name:     fmt.Sprintf("%s Stats", g.Name),
		Type:     discordgo.ChannelTypeGuildCategory,
		Topic:    "Here be stats",
		Position: 0,
	}
	cat, err := s.GuildChannelCreateComplex(g.ID, cat_data)
	if err != nil {
		log.Panic(err)
	}

	ch_data := discordgo.GuildChannelCreateData{
		Name:     fmt.Sprintf("%s Members: %d", g.Name, g.MemberCount),
		Type:     discordgo.ChannelTypeGuildVoice,
		Topic:    "eyyyyy",
		ParentID: cat.ID,
	}

	_, err = s.GuildChannelCreateComplex(g.ID, ch_data)
	if err != nil {
		log.Panic(err)
	}

	// ch_data = discordgo.GuildChannelCreateData{
	// 	Name:     fmt.Sprintf("test: %d", testseconds),
	// 	Type:     discordgo.ChannelTypeGuildVoice,
	// 	Topic:    "eyyyyy",
	// 	ParentID: cat.ID,
	// }

	// _, err = s.GuildChannelCreateComplex(g.ID, ch_data)
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func UpdateGuild(s *discordgo.Session, g *discordgo.Guild) {
	log.Printf("Updating guild %s\n", g.Name)
	parent_id, err := GetStatsCategoryID(s, g)
	if err != nil {
		log.Fatalln(err)
	}

	for _, ch := range g.Channels {
		if ch.ParentID == parent_id {
			if strings.HasPrefix(ch.Name, fmt.Sprintf("%s Members:", g.Name)) {
				s.ChannelEdit(ch.ID, &discordgo.ChannelEdit{
					Name: fmt.Sprintf("%s Members: %d", g.Name, g.MemberCount),
				})
			}
			// } else if strings.HasPrefix(ch.Name, "timetospoop:") {
			// 	log.Printf("Updating test channel with %d", testseconds)
			// 	st, err := s.ChannelEdit(ch.ID, &discordgo.ChannelEdit{
			// 		Name: fmt.Sprintf("test: %d", testseconds),
			// 	})
			// 	if err != nil {
			// 		log.Println(err)
			// 	} else {
			// 		log.Println(st.Name)
			// 	}
			// 	fmt.Println("did i get here")
			// }
		}
	}
}
