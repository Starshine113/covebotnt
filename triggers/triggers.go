package triggers

import "github.com/Starshine113/covebotnt/crouter"

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
		Aliases: []string{"Delete"},

		Description: "Remove a trigger",
		Usage:       "<trigger>",

		Permissions: crouter.PermLevelMod,
		Command:     remove,
	})
}
