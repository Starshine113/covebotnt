package ownercommands

import "github.com/starshine-sys/covebotnt/crouter"

// Init adds all the commands from this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "AdminDM",
		Description: "Send any user sharing a server with the bot a message, including attachment",
		Usage:       "<user ID> <message>",
		Permissions: crouter.PermLevelOwner,
		Command:     AdminDM,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Guilds",
		Description: "Show all guilds",
		Usage:       "",
		Permissions: crouter.PermLevelOwner,
		Command:     Guilds,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Guild",
		Description: "Get info about a specific guild",
		Usage:       "<id>",
		Permissions: crouter.PermLevelOwner,
		Command:     Guild,
	})

	router.AddCommand(&crouter.Command{
		Name:        "GetNewID",
		Aliases:     []string{"GetNewSnowflake", "GetSnowflake"},
		Description: "Get a new CoveBot snowflake",
		Permissions: crouter.PermLevelNone,
		Command:     snowflake,
	})

	router.AddCommand(&crouter.Command{
		Name:        "FixSnowflakes",
		Description: "Fix all mod log entries without a snowflake",
		Permissions: crouter.PermLevelOwner,
		Command:     fixSnowflakes,
	})
}
