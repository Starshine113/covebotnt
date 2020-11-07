package main

import "github.com/Starshine113/covebotnt/crouter"

func commandTree(ctx *crouter.Ctx) {
	guildSettings := globalSettings[ctx.Message.GuildID]
	ctx.AdditionalParams = map[string]interface{}{"config": config, "botVer": botVersion, "gitVer": string(gitOut), "startTime": startTime}
	err := router.Execute(ctx, &guildSettings, config.Bot.BotOwners)

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
