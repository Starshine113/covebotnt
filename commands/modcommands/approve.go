package modcommands

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/starshine-sys/covebotnt/etc"
	"github.com/starshine-sys/covebotnt/structs"
)

// GkApprove approves the given member, giving them the member roles
func GkApprove(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if ctx.RawArgs == "" {
		_, err = ctx.Send(crouter.ErrorEmoji + " No user supplied.")
		return err
	}
	member, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	guildConf := ctx.AdditionalParams["guildSettings"].(*structs.GuildSettings)
	for _, role := range guildConf.Gatekeeper.MemberRoles {
		ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, member.User.ID, role)
	}
	for _, role := range guildConf.Gatekeeper.GatekeeperRoles {
		ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, member.User.ID, role)
	}

	if guildConf.Gatekeeper.WelcomeChannel == "" {
		_, err = ctx.Send(fmt.Sprintf("%v No welcome channel set. Set one with `%vgatekeeper welcome-channel <channel>`.", crouter.WarnEmoji, ctx.GuildPrefix))
		if err != nil {
			return err
		}
	} else {
		var msgB bytes.Buffer
		tmpl, err := template.New("welcome").Parse(guildConf.Gatekeeper.WelcomeMessage)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(&msgB, member); err != nil {
			return err
		}
		_, err = ctx.Session.ChannelMessageSend(guildConf.Gatekeeper.WelcomeChannel, msgB.String())
		if err != nil {
			return err
		}
	}

	_, err = ctx.Send(fmt.Sprintf("%v **%v** approved **%v**.", crouter.SuccessEmoji, ctx.Author.Mention(), member.User.Mention()))
	if err != nil {
		return err
	}

	if guildConf.Moderation.ModLog == "" {
		_, err = ctx.Send(fmt.Sprintf("%v No mod log channel set. Set one with `%vmodlog <channel>`.", crouter.WarnEmoji, ctx.GuildPrefix))
		return
	}

	joined, err := member.JoinedAt.Parse()
	if err != nil {
		joined = time.Now().UTC()
	}

	guild, _ := ctx.Session.State.Guild(ctx.Message.GuildID)
	guildIcon := ""
	if guild != nil {
		guildIcon = guild.IconURL()
	}

	logEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "User Approved",
			IconURL: guildIcon,
		},
		Color:     0x0154C6,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Moderator",
				Value: fmt.Sprintf("**%s** (%v)", ctx.Author, ctx.Author.Mention()),
			},
			{
				Name:  "Target",
				Value: fmt.Sprintf("**%s** (%v)", member.User, member.User.Mention()),
			},
			{
				Name:  "Target joined at",
				Value: fmt.Sprintf("%v UTC (%v)", joined.Format("Jan 02 2006, 15:04:05"), etc.HumanizeTime(etc.DurationPrecisionMinutes, joined)),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Mod: %s (%v)", ctx.Author, ctx.Author.ID),
			IconURL: ctx.Author.AvatarURL("128"),
		},
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(guildConf.Moderation.ModLog, logEmbed)
	return
}
