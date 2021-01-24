package yagimport

import (
	"github.com/starshine-sys/covebotnt/bot"
	"github.com/starshine-sys/covebotnt/crouter"
)

type yag struct {
	*bot.Bot
}

// Init ...
func Init(router *crouter.Router) {
	y := &yag{router.Bot}

	y.Session.AddHandler(y.messageCreate)

	g := router.AddGroup(&crouter.Group{
		Name:        "YagImport",
		Description: "Manage importing moderation logs from YAGPDB.xyz",
		Command: &crouter.Command{
			Name:        "Settings",
			Description: "Show the server's current settings",
			Permissions: crouter.PermLevelMod,
			Command:     y.settings,
		},
	})

	g.AddCommand(&crouter.Command{
		Name:    "Toggle",
		Aliases: []string{"Enable", "Disable"},

		Description: "Toggle automatically adding to the moderation log when a user is warned through YAGPDB.xyz",

		Permissions: crouter.PermLevelMod,
		Command:     y.toggle,
	})

	g.AddCommand(&crouter.Command{
		Name: "Channel",

		Description: "Set the channel being listened to for moderation actions",

		Permissions: crouter.PermLevelMod,
		Command:     y.channel,
	})

	g.AddCommand(&crouter.Command{
		Name: "Import",

		Description: "Bulk import entries from the log channel",

		Permissions: crouter.PermLevelAdmin,
		Command:     y.bulk,
	})
}
