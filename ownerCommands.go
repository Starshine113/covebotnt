package main

import (
	"fmt"
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
)

func commandSetStatus(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	// this command needs at least 2 arguments
	if len(ctx.Args) < 2 {
		ctx.CommandError(&cbctx.ErrorNotEnoughArgs{2, len(ctx.Args)})
		return nil
	}

	// check the first arg -- boolean between -replace and -append
	if ctx.Args[0] == "-replace" {
		config.Bot.CustomStatus.Override = true
	} else if ctx.Args[0] == "-append" {
		config.Bot.CustomStatus.Override = false
	} else {
		ctx.CommandError(&cbctx.ErrorMissingRequiredArgs{"<-replace/-append> [-clear] <status string>", "<-replace/-append>"})
		return nil
	}

	// set custom status to the specified string
	config.Bot.CustomStatus.Status = strings.Join(ctx.Args[1:], " ")

	// check second argument -- if it's `-clear` the custom status is cleared
	if ctx.Args[1] == "-clear" {
		config.Bot.CustomStatus.Status = ""
	}

	// set the status
	newStatus := config.Bot.Prefixes[0] + "help | in " + fmt.Sprint(len(ctx.Session.State.Guilds)) + " guilds"
	if config.Bot.CustomStatus.Status != "" {
		newStatus += " | " + config.Bot.CustomStatus.Status
	}
	if config.Bot.CustomStatus.Override {
		newStatus = config.Bot.CustomStatus.Status
	}
	err = dg.UpdateStatus(0, newStatus)
	if err != nil {
		ctx.CommandError(err)
		return err
	}
	_, err = ctx.Send("Set status to `" + newStatus + "`")
	if err != nil {
		sugar.Errorf("Error when sending message: ", err)
		return err
	}
	sugar.Infof("Updated status to \"%v\"", newStatus)
	return nil
}
