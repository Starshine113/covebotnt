package commands

import (
	"fmt"
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// Nickname changes the bot's nickname
func Nickname(ctx *cbctx.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		ctx.CommandError(&cbctx.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0,
		})
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
