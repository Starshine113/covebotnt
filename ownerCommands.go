package main

import (
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/starshine-sys/covebotnt/crouter"
)

func commandSetStatus(ctx *crouter.Ctx) (err error) {
	if err = ctx.TriggerTyping(); err != nil {
		return err
	}

	// this command needs at least 1 arguments
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	// check the first arg
	switch ctx.Args[0] {
	case "listening":
		config.Bot.CustomStatus.Type = "listening"
	case "playing":
		config.Bot.CustomStatus.Type = "playing"
	case "watching":
		config.Bot.CustomStatus.Type = "watching"
	default:
		fallthrough
	case "clear":
		config.Bot.CustomStatus.Status = ""
		config.Bot.CustomStatus.Type = ""
	}

	if config.Bot.CustomStatus.Type == "" {
		err = dg.UpdateGameStatus(0, "")
		if err != nil {
			ctx.CommandError(err)
			return err
		}
		return nil
	}

	if len(ctx.Args) < 2 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 2,
			SuppliedArgs:    len(ctx.Args),
		})
		return nil
	}

	// set custom status to the specified string
	config.Bot.CustomStatus.Status = strings.Join(ctx.Args[1:], " ")

	err = updateStatus(ctx.Session)
	if err != nil {
		ctx.CommandError(err)
		sugar.Errorf("Error setting status: %v", err)
		return nil
	}

	_, err = ctx.Send("Set status to `" + config.Bot.CustomStatus.Status + "`")
	if err != nil {
		sugar.Errorf("Error when sending message: %v", err)
		return err
	}
	sugar.Infof("Updated status to \"%v\"", config.Bot.CustomStatus.Status)

	b, err := toml.Marshal(config)
	if err != nil {
		sugar.Errorf("Error marshaling toml config: %v", err)
		return
	}
	err = ioutil.WriteFile("config.toml", b, 0644)
	if err != nil {
		sugar.Errorf("Error writing config: %v", err)
		return
	}

	return nil
}
