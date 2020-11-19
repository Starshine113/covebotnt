package modcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// Warn warns the specified member
func Warn(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 2 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 2,
			SuppliedArgs:    len(ctx.Args),
		})
		return nil
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	warnReason := strings.Join(ctx.Args[1:], " ")

	dmChannel, err := ctx.Session.UserChannelCreate(member.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	warnMessage := fmt.Sprintf("%v You were warned in %v.\n**Reason:** %v", crouter.WarnEmoji, guild.Name, warnReason)
	_, err = ctx.Session.ChannelMessageSend(dmChannel.ID, warnMessage)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	currentLogs, err := ctx.Database.GetModLogs(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	var warnCount int
	for _, entry := range currentLogs {
		if entry.Type == "warn" {
			warnCount++
		}
	}

	var warnCountStr string
	switch fmt.Sprint(warnCount + 1)[len(fmt.Sprint(warnCount+1))-1] {
	case byte('1'):
		warnCountStr = fmt.Sprintf("%vst", warnCount+1)
	case byte('2'):
		warnCountStr = fmt.Sprintf("%vnd", warnCount+1)
	case byte('3'):
		warnCountStr = fmt.Sprintf("%vrd", warnCount+1)
	default:
		warnCountStr = fmt.Sprintf("%vth", warnCount+1)
	}

	_, err = ctx.Send(fmt.Sprintf("%v Warned **%v**, this is their %v warning.", crouter.SuccessEmoji, member.User.String(), warnCountStr))
	if err != nil {
		return err
	}

	entry, err := ctx.Database.AddToModLog(&cbdb.ModLogEntry{
		GuildID: ctx.Message.GuildID,
		UserID:  member.User.ID,
		ModID:   ctx.Author.ID,
		Type:    "warn",
		Reason:  warnReason,
	})
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	logEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.Author.String(),
			IconURL: ctx.Author.AvatarURL("256"),
		},
		Color:       0xe5da00,
		Title:       fmt.Sprintf("User warned (case #%v)", entry.ID),
		Description: fmt.Sprintf("**%v** (ID: %v)", member.User.String(), member.User.ID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Reason",
				Value:  warnReason,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Moderator ID: " + ctx.Author.ID,
		},
		Timestamp: entry.Time.Format(time.RFC3339),
	}

	modLog := ctx.AdditionalParams["guildSettings"].(*structs.GuildSettings).Moderation.ModLog

	if modLog == "" {
		_, err = ctx.Send(fmt.Sprintf("%v No mod log channel set. Set one with `%vmodlog <channel>`.", crouter.WarnEmoji, ctx.GuildPrefix))
		return
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(modLog, logEmbed)
	return
}
