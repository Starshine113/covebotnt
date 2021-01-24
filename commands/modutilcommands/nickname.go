package modutilcommands

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Nickname changes the bot's nickname
func Nickname(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		ctx.CommandError(err)
		return nil
	}
	nick := strings.Join(ctx.Args, " ")
	if len(nick) > 32 {
		nick = nick[:32]
	}
	err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, "@me", nick)
	if err != nil {
		ctx.CommandError(err)
		return err
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Bot nickname changed",
		Description: fmt.Sprintf("Successfully changed the bot's nickname to **%v**.", nick),
	})
	return
}
