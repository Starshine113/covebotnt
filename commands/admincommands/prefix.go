package admincommands

import (
	"context"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
)

func prefixShow(ctx *crouter.Ctx) (err error) {
	title := "Prefixes"
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == nil {
		title += " for " + guild.Name
	}

	x := make([]string, 0)
	for i, p := range ctx.Prefixes {
		if i == 0 {
			continue
		}
		if strings.HasPrefix(p, "<@") && strings.HasSuffix(p, ">") {
			x = append(x, p)
		} else {
			x = append(x, "`"+p+"`")
		}
	}

	_, err = ctx.Embed(title, strings.Join(x, "\n"), 0)
	return err
}

func prefixAdd(ctx *crouter.Ctx) (err error) {
	prefix := strings.Join(ctx.Args, " ")
	var matched bool

	err = ctx.Database.Pool.QueryRow(context.Background(), "select $1 = any(g.prefixes) from (select * from guild_settings where guild_id = $2) as g", prefix, ctx.Message.GuildID).Scan(&matched)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if matched {
		_, err = ctx.SendfNoAddXHandler("%v `%v` is already a prefix.", crouter.ErrorEmoji, prefix)
		return
	}

	if len(ctx.GuildSettings.Prefixes) > 22 {
		_, err = ctx.SendNoAddXHandler("This server has too many prefixes. Remove one to add a new prefix.")
		return
	}

	_, err = ctx.Database.Pool.Exec(context.Background(), "update public.guild_settings set prefixes = array_append(prefixes, $1) where guild_id = $2", prefix, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	if err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID); err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	_, err = ctx.SendfNoAddXHandler("%v Added prefix `%v`.", crouter.SuccessEmoji, prefix)
	return
}

func prefixRemove(ctx *crouter.Ctx) (err error) {
	prefix := strings.Join(ctx.Args, " ")
	var matched bool

	err = ctx.Database.Pool.QueryRow(context.Background(), "select $1 = any(g.prefixes) from (select * from guild_settings where guild_id = $2) as g", prefix, ctx.Message.GuildID).Scan(&matched)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	if !matched {
		_, err = ctx.SendfNoAddXHandler("%v `%v` is not a prefix.", crouter.ErrorEmoji, prefix)
		return
	}

	_, err = ctx.Database.Pool.Exec(context.Background(), "update public.guild_settings set prefixes = array_remove(prefixes, $1) where guild_id = $2", prefix, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	if err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID); err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	_, err = ctx.SendfNoAddXHandler("%v Removed prefix `%v`.", crouter.SuccessEmoji, prefix)
	return
}
