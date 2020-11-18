package starboard

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func (sb *Sb) createOrEditMessage(s *discordgo.Session, message *discordgo.Message, guildID string, count int, emoji *discordgo.Emoji) (err error) {
	embed, err := starboardEmbed(s, message, guildID)
	if err != nil {
		return err
	}
	msgContent := "**" + fmt.Sprint(count) + "** " + emoji.MessageFormat() + " <#" + message.ChannelID + ">"

	gs, err := sb.Pool.GetGuildSettings(guildID)
	if err != nil {
		return err
	}

	if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content: &msgContent,
			Embed:   embed,
			ID:      m,
			Channel: gs.Starboard.StarboardID,
		})
		if err != nil {
			return err
		}
	} else {
		starboardMsg, err := s.ChannelMessageSendComplex(gs.Starboard.StarboardID, &discordgo.MessageSend{
			Content: msgContent,
			Embed:   embed,
		})
		if err != nil {
			return err
		}
		err = sb.Pool.InsertStarboardEntry(message.ID, message.ChannelID, guildID, starboardMsg.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sb *Sb) deleteStarboardMessage(s *discordgo.Session, starboardMessage, guildID string) error {
	gs, err := sb.Pool.GetGuildSettings(guildID)
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(gs.Starboard.StarboardID, starboardMessage)
	if err != nil {
		return err
	}
	err = sb.Pool.DeleteStarboardEntry(starboardMessage)
	return err
}

func starboardEmbed(s *discordgo.Session, message *discordgo.Message, guildID string) (embed *discordgo.MessageEmbed, err error) {
	name := message.Author.Username
	if message.WebhookID == "" {
		member, err := s.State.Member(guildID, message.Author.ID)
		if err != nil {
			member, err = s.GuildMember(guildID, message.Author.ID)
		}
		if err == nil && member.Nick != "" {
			name = member.Nick
		}
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
			{Name: "Source", Value: "[Jump to message](https://discordapp.com/channels/" + guildID + "/" + message.ChannelID + "/" + message.ID + ")", Inline: false},
		},
		Timestamp: string(message.Timestamp),
		Color:     0xede21e,
		Image: &discordgo.MessageEmbedImage{
			URL: attachmentLink,
		},
	}
	return embed, nil
}
