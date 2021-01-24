package ownercommands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/starshine-sys/covebotnt/commands/usercommands"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Guilds ...
func Guilds(ctx *crouter.Ctx) (err error) {
	guilds := make([]string, 0)
	for _, g := range ctx.Session.State.Guilds {
		guilds = append(guilds, fmt.Sprintf("%v (%v)", g.Name, g.ID))
	}

	if len(strings.Join(guilds, "\n")) < 2000 {
		_, err = ctx.Embedf(fmt.Sprintf("Guilds (%v)", len(ctx.Session.State.Guilds)), "```%v```", strings.Join(guilds, "\n"))
		return
	}

	reader := bytes.NewReader([]byte(strings.Join(guilds, "\n")))

	file := discordgo.File{
		Name:   fmt.Sprintf("guilds-%v.txt", time.Now().UTC().Format("2006-01-02-15-04-05")),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})
	return
}

// Guild ...
func Guild(ctx *crouter.Ctx) (err error) {
	if ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	guild, err := ctx.Session.State.Guild(ctx.Args[0])
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	msg, err := ctx.Send("Working, please wait...")
	if err != nil {
		return err
	}

	guildCreated, err := discordgo.SnowflakeTimestamp(guild.ID)
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
				Value:  fmt.Sprintf("%v\n(%v ago)", guildCreated.Format("Jan _2 2006, 15:04:05 MST"), usercommands.PrettyDurationString(time.Since(guildCreated))),
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
