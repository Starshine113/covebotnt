package modcommands

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/structs"

	"github.com/bwmarrin/discordgo"
)

var idRegex = regexp.MustCompile("\\d{16,}")

// Ban bans the specified user from the server
func Ban(ctx *crouter.Ctx) (err error) {
	reason := "None"

	if len(ctx.Args) == 0 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<user ID>",
			MissingArgs:  "<user ID>",
		})
		return nil
	}

	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	var user *discordgo.User
	var m bool
	member, err := ctx.ParseMember(ctx.Args[0])
	if err == nil {
		m = true
		user = member.User
	} else {
		if !idRegex.MatchString(ctx.Args[0]) {
			_, err = ctx.Sendf("`%v` is not a valid Discord user ID.", ctx.Args[0])
			return err
		}
		m = false
		user, err = ctx.Session.User(ctx.Args[0])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	var guild *discordgo.Guild
	if m {
		var memberRoles, modRoles discordgo.Roles

		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}

		mod, err := ctx.ParseMember(ctx.Author.ID)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}

		for _, r := range guild.Roles {
			for _, m := range mod.Roles {
				if r.ID == m {
					modRoles = append(modRoles, r)
				}
			}
			for _, m := range member.Roles {
				if r.ID == m {
					memberRoles = append(memberRoles, r)
				}
			}
		}

		sort.Sort(modRoles)
		sort.Sort(memberRoles)

		if len(modRoles) == 0 {
			if guild.OwnerID != mod.User.ID {
				_, err = ctx.Send("You're not high enough in the role hierarchy to do that.")
				return err
			}
		}
		if len(memberRoles) != 0 {
			if modRoles[0].Position <= memberRoles[0].Position {
				_, err = ctx.Send("You're not high enough in the role hierarchy to do that.")
				return err
			}
		}
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
			_, err = ctx.Send(fmt.Sprintf("%v User **%v** is already banned.\n(Reason: __%v__)", crouter.ErrorEmoji, user.String(), reason))
			return
		}
	}

	formattedReason := fmt.Sprintf("%v: %v", ctx.Author.String(), reason)

	if m && !user.Bot {
		dmChannel, err := ctx.Session.UserChannelCreate(member.User.ID)
		if err != nil {
			ctx.CommandError(err)
			return nil
		}

		_, err = ctx.Session.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("You were banned from %v. Reason: %v", guild.Name, reason))
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, user.ID, formattedReason, 2)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Send(fmt.Sprintf("%v Banned **%v** with reason: __%v__", crouter.SuccessEmoji, user.String(), reason))
	if err != nil {
		return err
	}

	entry, err := ctx.Database.AddToModLog(&cbdb.ModLogEntry{
		GuildID: ctx.Message.GuildID,
		UserID:  user.ID,
		ModID:   ctx.Author.ID,
		Type:    "ban",
		Reason:  reason,
		Time:    time.Now().UTC(),
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
		_, err = ctx.Send(fmt.Sprintf("%v No mod log channel set. Set one with `%vmodlog <channel>`.", crouter.WarnEmoji, ctx.GuildPrefix))
		return
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(modLog, logEmbed)
	return
}
