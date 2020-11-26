package usercommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// GuildInfo command shows information about the current guild
func GuildInfo(ctx *crouter.Ctx) (err error) {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	msg, err := ctx.Send("Working, please wait...")
	if err != nil {
		return err
	}

	guildCreated, err := discordgo.SnowflakeTimestamp(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	guildCreated = guildCreated.UTC()

	var categoryCount, totalChanCount, textChanCount, voiceChanCount int
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			categoryCount++
		} else {
			totalChanCount++
		}
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChanCount++
		}
		if channel.Type != discordgo.ChannelTypeGuildCategory && channel.Type != discordgo.ChannelTypeGuildVoice {
			textChanCount++
		}
	}

	features := guild.Features
	if len(features) == 0 {
		features = []string{"NONE"}
	}

	ownerString := "<@" + guild.OwnerID + ">"
	owner, err := ctx.Session.User(guild.OwnerID)
	if err == nil {
		ownerString = owner.String()
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Info for %v", guild.Name),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ID: " + guild.ID + " | Created ",
		},
		Color:     0x7289da,
		Timestamp: guildCreated.Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL(),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Owner",
				Value:  ownerString,
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprint(guild.MemberCount),
				Inline: true,
			},
			{
				Name:   "Level",
				Value:  fmt.Sprintf("%v (%v boosts)", guild.PremiumTier, guild.PremiumSubscriptionCount),
				Inline: true,
			},
			{
				Name:   "Emoji",
				Value:  fmt.Sprint(len(guild.Emojis)),
				Inline: true,
			},
			{
				Name:   "Prefix",
				Value:  "`" + ctx.GuildPrefix + "`",
				Inline: true,
			},
			{
				Name:   "Roles",
				Value:  fmt.Sprintf("%v", len(guild.Roles)),
				Inline: true,
			},
			{
				Name:   "Channels",
				Value:  fmt.Sprintf("%v (in %v categories)\n<:textchannel:770274583223336990> %v | <:voicechannel:770274509012992020> %v", totalChanCount, categoryCount, textChanCount, voiceChanCount),
				Inline: true,
			},
			{
				Name:   "Created at",
				Value:  fmt.Sprintf("%v\n(%v ago)", guildCreated.Format("Jan _2 2006, 15:04:05 MST"), PrettyDurationString(time.Since(guildCreated))),
				Inline: true,
			},
			{
				Name:   "Features",
				Value:  "```" + strings.Join(features, ", ") + "```",
				Inline: false,
			},
		},
	}

	content := ""
	_, err = ctx.Edit(msg, &discordgo.MessageEdit{
		Content: &content,
		Embed:   embed,
	})
	return
}
