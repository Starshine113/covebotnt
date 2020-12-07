package modutilcommands

import (
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func invite(ctx *crouter.Ctx) (err error) {
	channel := ctx.Channel
	if len(ctx.Args) > 0 {
		channel, err = ctx.ParseChannel(strings.Join(ctx.Args, " "))
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
	}

	invite, err := ctx.Session.ChannelInviteCreate(channel.ID, discordgo.Invite{
		MaxAge:    0,
		MaxUses:   0,
		Temporary: false,
		Unique:    true,
	})
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.SendfNoAddXHandler("Created invite **%v** for channel %v.\nLink:\n<https://discord.gg/%v>", invite.Code, invite.Channel.Mention(), invite.Code)
	return
}
