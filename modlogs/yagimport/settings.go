package yagimport

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/starshine-sys/covebotnt/crouter"
)

func (y *yag) settings(ctx *crouter.Ctx) (err error) {
	gs, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var c string
	if gs.YAG.Channel == "" {
		c = "None"
	} else {
		c = "<#" + gs.YAG.Channel + ">"
	}

	_, err = ctx.Embedf("YAGPDB.xyz import", "- Enabled: `%v`\n- Channel: %v", gs.YAG.Enabled, c)
	return
}

func (y *yag) toggle(ctx *crouter.Ctx) (err error) {
	gs, err := ctx.Database.GetGuildSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	b := !gs.YAG.Enabled

	if len(ctx.Args) > 0 {
		v, err := strconv.ParseBool(ctx.Args[0])
		if err == nil {
			b = v
		}
	}

	commandTag, err := ctx.Database.Pool.Exec(context.Background(), "update public.yag_import set enabled = $1 where guild_id = $2", b, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		_, err = ctx.CommandError(errors.New("no rows affected"))
		return err
	}

	ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)

	if b {
		_, err = ctx.Send("Enabled automatic logging of new mod log entries.")
	} else {
		_, err = ctx.Send("Disabled automatic logging of new mod log entries.")
	}
	return
}

func (y *yag) channel(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	channel, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if channel.GuildID != ctx.Message.GuildID {
		_, err = ctx.SendfNoAddXHandler("Channel %v is not in this server.", channel.Mention())
		return err
	}

	commandTag, err := ctx.Database.Pool.Exec(context.Background(), "update public.yag_import set log_channel = $1 where guild_id = $2", channel.ID, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		_, err = ctx.CommandError(errors.New("no rows affected"))
		return err
	}

	ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)

	_, err = ctx.Sendf("Set the YAGPDB.xyz log channel to %v.", channel.Mention())
	return
}
