package main

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func createOrEditMessage(s *discordgo.Session, message *discordgo.Message, guildID string, count int, emoji *discordgo.Emoji) (err error) {
	embed, err := starboardEmbed(message, guildID)
	if err != nil {
		return err
	}
	msgContent := "**" + fmt.Sprint(count) + "** " + emoji.MessageFormat() + " <#" + message.ChannelID + ">"

	if val, ok := messageIDMap[message.ID]; ok {
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content: &msgContent,
			Embed:   embed,
			ID:      val,
			Channel: globalSettings[guildID].Starboard.StarboardID,
		})
		if err != nil {
			return err
		}
	} else {
		starboardMsg, err := s.ChannelMessageSendComplex(globalSettings[guildID].Starboard.StarboardID, &discordgo.MessageSend{
			Content: msgContent,
			Embed:   embed,
		})
		if err != nil {
			return err
		}
		err = insertStarboardEntry(message.ID, message.ChannelID, guildID, starboardMsg.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteStarboardMessage(s *discordgo.Session, starboardMessage, guildID string) error {
	err := s.ChannelMessageDelete(globalSettings[guildID].Starboard.StarboardID, starboardMessage)
	if err != nil {
		return err
	}
	err = deleteStarboardEntry(starboardMessage)
	return err
}

func starboardEmbed(message *discordgo.Message, guildID string) (embed *discordgo.MessageEmbed, err error) {
	name := message.Author.Username
	if message.Member != nil {
		name = message.Member.Nick
	}
	var attachmentLink string
	if len(message.Attachments) > 0 {
		match, _ := regexp.MatchString("\\.(png|jpg|jpeg|gif|webp)$", message.Attachments[0].URL)
		if match {
			attachmentLink = message.Attachments[0].URL
		}
	}

	embed = &discordgo.MessageEmbed{
		Description: message.Content,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: message.Author.AvatarURL("256"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ID: " + message.ID,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Source", Value: "[Jump!](https://discordapp.com/channels/" + guildID + "/" + message.ChannelID + "/" + message.ID + ")", Inline: false},
		},
		Timestamp: string(message.Timestamp),
		Color:     0xede21e,
		Image: &discordgo.MessageEmbedImage{
			URL: attachmentLink,
		},
	}
	return embed, nil
}
