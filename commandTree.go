package main

import (
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/commands"
)

var allCommands string = `
All commands
============
|-- ping
|-- help
|-- setstatus <-replace|-append> <status string|-clear>
|-- starboard
|   |-- channel <channel>
|   |-- limit <int>
|-- modroles
|-- echo [-channel <channel>] <message string>
|-- prefix [prefix string]
|-- steal <emoji url/emoji> [name string]
|-- enlarge <emoji>
|-- commands
|-- setnote <user> <note>
|-- notes <user>
|-- delnote <id>`

const successEmoji, errorEmoji string = "✅", "❌"

func commandTree(ctx *cbctx.Ctx) {
	var err error
	if err != nil {
		sugar.Errorf("Error getting context: %v", err)
	}

	switch ctx.Command {
	case "ping":
		err = commands.Ping(ctx)
	case "help":
		ctx.AdditionalParams = map[string]interface{}{"config": config}
		err = commands.Help(ctx)
	case "setstatus":
		err = commandSetStatus(ctx)
	case "starboard":
		err = commandStarboard(ctx)
	case "modroles":
		err = commandModRoles(ctx.Args, ctx.Session, ctx.Message)
	case "echo":
		err = commandEcho(ctx)
	case "prefix":
		err = commandPrefix(ctx.Args, ctx.Session, ctx.Message)
	case "commands":
		_, err = ctx.Send(fmt.Sprintf("```%v```", allCommands))
	case "steal":
		err = commandSteal(ctx)
	case "enlarge":
		err = commandEnlarge(ctx)
	case "notes":
		err = commandNotes(ctx)
	case "setnote", "addnote":
		err = commandSetNote(ctx)
	case "delnote", "removenote":
		err = commandDelNote(ctx)
	}

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
