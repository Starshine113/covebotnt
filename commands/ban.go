package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/structs"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// Ban bans the specified user from the server
func Ban(ctx *cbctx.Ctx) (err error) {
	reason := "None"

	if len(ctx.Args) == 0 {
		ctx.CommandError(&cbctx.ErrorMissingRequiredArgs{
			RequiredArgs: "<user ID>",
			MissingArgs:  "<user ID>",
		})
		return nil
	}

	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	user, err := ctx.Session.User(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	currentBans, err := ctx.Session.GuildBans(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	for _, ban := range currentBans {
		if ban.User.ID == user.ID {
			reason = "None"
			if ban.Reason != "" {
				reason = ban.Reason
			}
			_, err = ctx.Send(fmt.Sprintf("%v User **%v** is already banned.\n(Reason: __%v__)", cbctx.ErrorEmoji, user.String(), reason))
			return
		}
	}

	formattedReason := fmt.Sprintf("%v: %v", ctx.Author.String(), reason)

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, user.ID, formattedReason, 2)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Send(fmt.Sprintf("%v Banned **%v** with reason: __%v__", cbctx.SuccessEmoji, user.String(), reason))
	if err != nil {
		return err
	}

	entry, err := ctx.Database.AddToModLog(&cbdb.ModLogEntry{
		GuildID: ctx.Message.GuildID,
		UserID:  user.ID,
		ModID:   ctx.Author.ID,
		Type:    "ban",
		Reason:  reason,
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
		Color:       0xc13030,
		Title:       fmt.Sprintf("User banned (case #%v)", entry.ID),
		Description: fmt.Sprintf("**%v** (ID: %v)", user.String(), user.ID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Reason",
				Value:  reason,
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
		_, err = ctx.Send(fmt.Sprintf("%v No mod log channel set. Set one with `%vmodlog <channel>`.", cbctx.WarnEmoji, ctx.GuildPrefix))
		return
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(modLog, logEmbed)
	return
}
