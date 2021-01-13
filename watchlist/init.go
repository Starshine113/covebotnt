package watchlist

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/wlog"
	"github.com/bwmarrin/discordgo"
)

type uh struct {
	pool  *cbdb.Db
	sugar *wlog.Wlog
}

// Init ...
func Init(r *crouter.Router, pool *cbdb.Db, sugar *wlog.Wlog) {
	u := &uh{
		pool:  pool,
		sugar: sugar,
	}

	r.Bot.Session.AddHandler(u.onJoin)

	p := r.AddGroup(&crouter.Group{
		Name:        "Watchlist",
		Description: "Manage the server's watchlist",
		Command: &crouter.Command{
			Name:        "Show",
			Description: "None",
			Permissions: crouter.PermLevelMod,
			Command:     watchlist,
		},
	})

	p.AddCommand(&crouter.Command{
		Name:        "Add",
		Description: "Add a user to the watchlist",
		Permissions: crouter.PermLevelMod,
		Command:     watchlistAdd,
	})

	p.AddCommand(&crouter.Command{
		Name:        "Remove",
		Description: "Remove a user from the watchlist",
		Permissions: crouter.PermLevelMod,
		Command:     watchlistRemove,
	})

	p.AddCommand(&crouter.Command{
		Name:        "Channel",
		Description: "Set the watchlist notification channel",
		Permissions: crouter.PermLevelMod,
		Command:     setChannel,
	})
}

func setChannel(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		_, err = ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<channel>",
			MissingArgs:  "<channel>",
		})
		return err
	}
	channel, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if channel.GuildID != ctx.Message.GuildID {
		ctx.CommandError(errors.New("channel must be in current guild"))
		return
	}
	err = ctx.Database.SetWatchlistChannel(ctx.Message.GuildID, channel.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}
	_, err = ctx.Send("Set the watchlist channel to " + channel.Mention() + ".")
	return
}

func watchlist(ctx *crouter.Ctx) (err error) {
	b, err := ctx.Database.GetWatchlist(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	var x string
	for _, c := range b {
		x += fmt.Sprintf("<@%v>\n", c)
	}
	if len(b) == 0 {
		x = "No channels are blacklisted."
	}
	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Channel blacklist",
		Description: x,
	})
	return err
}
