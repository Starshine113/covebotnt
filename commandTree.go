package main

import (
	"github.com/Starshine113/covebotnt/levels"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/commands"
)

const successEmoji, errorEmoji string = "✅", "❌"

func commandTree(ctx *cbctx.Ctx) {
	var err error
	if err != nil {
		sugar.Errorf("Error getting context: %v", err)
	}

	if ctx.Match("ping") {
		err = commands.Ping(ctx)
	} else if ctx.Match("help") {
		ctx.AdditionalParams = map[string]interface{}{"config": config, "botVer": botVersion, "gitVer": string(gitOut)}
		err = commands.Help(ctx)
	} else if ctx.Match("setstatus") {
		err = commandSetStatus(ctx)
	} else if ctx.Match("starboard", "sb") {
		err = commandStarboard(ctx)
	} else if ctx.Match("echo", "say") {
		err = commandEcho(ctx)
	} else if ctx.Match("steal", "addemote", "addemoji") {
		err = commandPrefix(ctx)
	} else if ctx.Match("enlarge", "emote", "emoji") {
		err = commandEnlarge(ctx)
	} else if ctx.Match("notes") {
		err = commandNotes(ctx)
	} else if ctx.Match("setnote", "addnote") {
		err = commandSetNote(ctx)
	} else if ctx.Match("delnote", "removenote") {
		err = commandDelNote(ctx)
	} else if ctx.Match("i", "info", "userinfo", "profile", "whois") {
		err = commands.UserInfo(ctx)
	} else if ctx.Match("si", "serverinfo", "guildinfo") {
		err = commands.GuildInfo(ctx)
	} else if ctx.Match("hello", "hi", "henlo", "heya", "heyo") {
		err = commands.Hello(ctx)
	} else if ctx.Match("export") {
		err = commandExport(ctx)
	} else if ctx.Match("archive") {
		return
		//err = commandArchive(ctx)
	} else if ctx.Match("level", "lvl", "rank") {
		err = levels.CommandLevel(ctx)
	} else if ctx.Match("leaderboard") {
		err = levels.CommandLeaderboard(ctx)
	} else if ctx.MatchRegex("^mod[-_]?roles?$") {
		err = commandModRoles(ctx)
	} else if ctx.MatchRegex("^helper[-_]?roles?$") {
		err = commandHelperRoles(ctx)
	}

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
