package modcommands

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// GkApprove approves the given member, giving them the member roles
func GkApprove(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<member>",
			MissingArgs:  "<member>",
		})
		return
	}
	member, err := ctx.ParseMember(strings.Join(ctx.Args, " "))
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

	logEmbed := &discordgo.MessageEmbed{
		Title:       "User approved",
		Description: fmt.Sprintf("%v (%v) was approved by %v (%v).", member.User.String(), member.User.ID, ctx.Author.String(), ctx.Author.ID),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Moderator ID: " + ctx.Author.ID,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(guildConf.Moderation.ModLog, logEmbed)
	return
}
