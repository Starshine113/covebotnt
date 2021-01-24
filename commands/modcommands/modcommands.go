package modcommands

import "github.com/starshine-sys/covebotnt/crouter"

// Init adds all the commands from this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "MajorityOf",
		Aliases:     []string{"Majority"},
		Description: "Get the majority of a number with abstains",
		Usage:       "<count> [abstains]",
		Permissions: crouter.PermLevelNone,
		Command:     MajorityOf,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Warn",
		Description: "Warn a user",
		Usage:       "<user> <reason>",
		Permissions: crouter.PermLevelHelper,
		Command:     Warn,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModLogs",
		Description: "Show a user's modlogs",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     ModLogs,
	})

	router.AddCommand(&crouter.Command{
		Name:        "LogMute",
		Description: "Log a user mute made with another bot",
		Usage:       "<user> [-d <duration>] [-hardmute] [reason]",
		Permissions: crouter.PermLevelHelper,
		Command:     LogMute,
	})

	router.AddCommand(&crouter.Command{
		Name:        "BGC",
		Aliases:     []string{"BackgroundCheck"},
		Description: "Quickly check a user",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     BGC,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Ban",
		Aliases:     []string{"HeartfeltConfession"},
		Description: "Ban users by ID",
		Usage:       "<user ID> [reason]",
		Permissions: crouter.PermLevelMod,
		Command:     Ban,
	})

	router.AddCommand(&crouter.Command{
		Name:        "MVC",
		Description: "Get the mvc majority",
		Usage:       "[abstains]",
		Permissions: crouter.PermLevelMod,
		Command:     MVC,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Export",
		Description: "Export this server's mod logs",
		Usage:       "",
		Permissions: crouter.PermLevelAdmin,
		Command:     Export,
	})
}
