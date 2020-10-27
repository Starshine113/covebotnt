package main

import (
	"github.com/Starshine113/covebotnt/cbctx"
)

const successEmoji, errorEmoji string = "✅", "❌"

func commandTree(ctx *cbctx.Ctx) {
	var err error
	if err != nil {
		sugar.Errorf("Error getting context: %v", err)
	}

	guildSettings := globalSettings[ctx.Message.GuildID]
	ctx.AdditionalParams = map[string]interface{}{"config": config, "botVer": botVersion, "gitVer": string(gitOut)}
	err = router.Execute(ctx, &guildSettings, config.Bot.BotOwners)

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
