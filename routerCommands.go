package main

import (
	"github.com/Starshine113/covebotnt/commands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/levels"
)

func addUserCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "ping",
		Description: "Ping pong!",
		Usage:       "ping",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Ping,
	})

	router.AddCommand(&crouter.Command{
		Name:        "enlarge",
		Aliases:     []string{"emote", "emoji"},
		Description: "Enlarges a custom emoji",
		Usage:       "enlarge <emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     commandEnlarge,
	})

	router.AddCommand(&crouter.Command{
		Name:        "about",
		Description: "Show some info about the bot",
		Usage:       "about",
		Permissions: crouter.PermLevelNone,
		Command:     commands.About,
	})

	router.AddCommand(&crouter.Command{
		Name:        "userinfo",
		Aliases:     []string{"i", "info", "whois", "profile"},
		Description: "Show information about a user (or yourself)",
		Usage:       "info [user]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.UserInfo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "serverinfo",
		Aliases:     []string{"si", "guildinfo"},
		Description: "Show information about the current server",
		Usage:       "serverinfo",
		Permissions: crouter.PermLevelNone,
		Command:     commands.GuildInfo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "hello",
		Aliases:     []string{"hi", "henlo", "heya", "heyo"},
		Description: "Say hi to the bot",
		Usage:       "hello",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Hello,
	})

	router.AddCommand(&crouter.Command{
		Name:        "lvl",
		Aliases:     []string{"level", "rank"},
		Description: "Show your (or another user's) level",
		Usage:       "lvl [user]",
		Permissions: crouter.PermLevelNone,
		Command:     levels.CommandLevel,
	})

	router.AddCommand(&crouter.Command{
		Name:        "leaderboard",
		Description: "Show the server's leaderboard",
		Usage:       "leaderboard [page]",
		Permissions: crouter.PermLevelNone,
		Command:     levels.CommandLeaderboard,
	})
}

func addHelperCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "notes",
		Description: "Shows a user's notes",
		Usage:       "notes <user>",
		Permissions: crouter.PermLevelHelper,
		Command:     commandNotes,
	})

	router.AddCommand(&crouter.Command{
		Name:        "setnote",
		Aliases:     []string{"addnote"},
		Description: "Set a note for a user",
		Usage:       "setnote <user> <note>",
		Permissions: crouter.PermLevelHelper,
		Command:     commandSetNote,
	})
}

func addModCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "delnote",
		Aliases:     []string{"removenote"},
		Description: "Remove a note by ID",
		Usage:       "delnote <id>",
		Permissions: crouter.PermLevelMod,
		Command:     commandDelNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "echo",
		Aliases:     []string{"say", "send"},
		Description: "Make the bot say something",
		Usage:       "echo [-ch <channel>] <message>",
		Permissions: crouter.PermLevelMod,
		Command:     commandEcho,
	})

	router.AddCommand(&crouter.Command{
		Name:        "steal",
		Aliases:     []string{"addemote", "addemoji"},
		Description: "Steal an emote by URL + name, or usage in message (with Nitro)",
		Usage:       "steal <emoji: url/emoji> [name]",
		Permissions: crouter.PermLevelMod,
		Command:     commandSteal,
	})

	router.AddCommand(&crouter.Command{
		Name:        "starboard",
		Aliases:     []string{"sb"},
		Description: "Manage the server's starboard",
		Usage:       "starboard <channel|limit|emoji>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboard,
	})
}

func addAdminCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "export",
		Description: "Export this server's notes",
		Usage:       "export",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandExport,
	})

	router.AddCommand(&crouter.Command{
		Name:        "modroles",
		Aliases:     []string{"mod-roles", "mod_roles", "mod-role", "mod_role", "modrole"},
		Description: "List/modify this server's mod roles",
		Usage:       "modroles [add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandModRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "helperroles",
		Aliases:     []string{"help-roles", "helper_roles", "helper-role", "helper_role", "helperrole"},
		Description: "List/modify this server's helper roles",
		Usage:       "helperroles [add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandHelperRoles,
	})
}

func addOwnerCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "setstatus",
		Description: "Set the bot's status",
		Usage:       "setstatus <-replace/-append> [<status>|-clear]",
		Permissions: crouter.PermLevelOwner,
		Command:     commandSetStatus,
	})
}