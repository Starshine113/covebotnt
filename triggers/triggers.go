package triggers

import "github.com/starshine-sys/covebotnt/crouter"

// Init adds this package's commands to the router
func Init(router *crouter.Router) {
	g := router.AddGroup(&crouter.Group{
		Name:        "Triggers",
		Aliases:     []string{"Trigger", "AutoResponse", "AutoResponses", "AR"},
		Description: "Manage triggers",
		Command: &crouter.Command{
			Name:        "Show",
			Description: "Show a list of this server's triggers, or a specific trigger",
			Usage:       "[trigger]",
			Permissions: crouter.PermLevelNone,
			GuildOnly:   true,
			Command:     show,
		},
	})

	g.AddCommand(&crouter.Command{
		Name:    "Add",
		Aliases: []string{"+", "Create"},

		Description:     "Add a trigger",
		LongDescription: "Add a trigger <trigger> with the response <response>.\n*The trigger and response must be separated by a newline.*",
		Usage:           "<trigger>\n<response>",

		Permissions: crouter.PermLevelMod,
		Command:     add,
	})

	g.AddCommand(&crouter.Command{
		Name:    "Remove",
		Aliases: []string{"Delete", "Yeet"},

		Description: "Remove a trigger",
		Usage:       "<trigger>",

		Permissions: crouter.PermLevelMod,
		Command:     remove,
	})

	g.AddCommand(&crouter.Command{
		Name: "Export",

		Description: "Exports all triggers as a .json file",

		Permissions: crouter.PermLevelMod,
		Command:     cmdexport,
	})

	g.AddCommand(&crouter.Command{
		Name: "Import",

		Description: "Imports triggers from a previous export",

		Permissions: crouter.PermLevelMod,
		Command:     cmdimport,
	})
}
