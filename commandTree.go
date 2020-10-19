package main

import (
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
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
|-- commands`

const successEmoji, errorEmoji string = "✅", "❌"

func commandTree(ctx *cbctx.Ctx) {
	var err error
	if err != nil {
		sugar.Errorf("Error getting context: %v", err)
	}

	switch ctx.Command {
	case "ping":
		err = commandPing(ctx)
	case "help":
		ctx.AdditionalParams = []interface{}{config}
		err = commandHelp(ctx)
	case "setstatus":
		err = commandSetStatus(ctx)
	case "starboard":
		err = commandStarboard(ctx)
	case "modroles":
		err = commandModRoles(ctx.Args, ctx.Session, ctx.Message)
	case "echo":
		err = commandEcho(ctx.Args, ctx.Session, ctx.Message)
	case "prefix":
		err = commandPrefix(ctx.Args, ctx.Session, ctx.Message)
	case "commands":
		_, err = ctx.Send(fmt.Sprintf("```%v```", allCommands))
	case "steal":
		err = commandSteal(ctx)
	case "enlarge":
		err = commandEnlarge(ctx)
	}

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
