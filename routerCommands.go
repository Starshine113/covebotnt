package main

import (
	"github.com/Starshine113/covebotnt/commands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/notes"
	"github.com/bwmarrin/discordgo"
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
		Aliases:     []string{"e", "emote", "emoji", "enlorge"},
		Description: "Enlarges a custom emoji",
		Usage:       "enlarge <emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Enlarge,
	})

	router.AddCommand(&crouter.Command{
		Name:        "color",
		Aliases:     []string{"colour"},
		Description: "Previews a color",
		Usage:       "color <color>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Color,
	})

	router.AddCommand(&crouter.Command{
		Name:        "avatar",
		Aliases:     []string{"pfp", "a"},
		Description: "Show a user's avatar",
		Usage:       "avatar [user]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Avatar,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "snowflake",
		Aliases:     []string{"idtime"},
		Description: "Get timestamps from the given ID(s)",
		Usage:       "snowflake [...IDs]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Snowflake,
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
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "roleinfo",
		Aliases:     []string{"ri"},
		Description: "Show information about a role",
		Usage:       "roleinfo <role>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.RoleInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "serverinfo",
		Aliases:     []string{"si", "guildinfo"},
		Description: "Show information about the current server",
		Usage:       "serverinfo",
		Permissions: crouter.PermLevelNone,
		Command:     commands.GuildInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "majorityof",
		Aliases:     []string{"majority"},
		Description: "Get the majority of a number with abstains",
		Usage:       "majorityof <count> [abstains]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.MajorityOf,
	})

	router.AddCommand(&crouter.Command{
		Name:        "hello",
		Aliases:     []string{"hi", "henlo", "heya", "heyo"},
		Description: "Say hi to the bot",
		Usage:       "hello",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Hello,
		GuildOnly:   true,
	})
}

func addHelperCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "notes",
		Description: "Shows a user's notes",
		Usage:       "notes <user>",
		Permissions: crouter.PermLevelHelper,
		Command:     notes.CommandNotes,
	})

	router.AddCommand(&crouter.Command{
		Name:        "setnote",
		Aliases:     []string{"addnote"},
		Description: "Set a note for a user",
		Usage:       "setnote <user> <note>",
		Permissions: crouter.PermLevelHelper,
		Command:     notes.CommandSetNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "warn",
		Description: "Warn a user",
		Usage:       "warn <user> <reason>",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.Warn,
	})

	router.AddCommand(&crouter.Command{
		Name:        "modlogs",
		Description: "Show a user's modlogs",
		Usage:       "warn <user>",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.ModLogs,
	})

	router.AddCommand(&crouter.Command{
		Name:        "logmute",
		Description: "Log a user mute made with another bot",
		Usage:       "logmute <user> [-d <duration>] [-hardmute] [reason]",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.LogMute,
	})
}

func addModCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "delnote",
		Aliases:     []string{"removenote"},
		Description: "Remove a note by ID",
		Usage:       "delnote <id>",
		Permissions: crouter.PermLevelMod,
		Command:     notes.CommandDelNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "nickname",
		Aliases:     []string{"nick"},
		Description: "Change the bot's nickname",
		Usage:       "nickname <new nickname>",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Nickname,
	})

	router.AddCommand(&crouter.Command{
		Name:        "echo",
		Aliases:     []string{"say", "send"},
		Description: "Make the bot say something",
		Usage:       "echo [-ch <channel>] <message>",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Echo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "steal",
		Aliases:     []string{"addemote", "addemoji"},
		Description: "Steal an emote by URL + name, or usage in message (with Nitro)",
		Usage:       "steal <emoji: url/emoji> [name]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Steal,
	})

	starboard := router.AddGroup(&crouter.Group{
		Name:        "starboard",
		Aliases:     []string{"sb"},
		Description: "Manage the server's starboard",
		Command: &crouter.Command{
			Name:        "show",
			Description: "Show the server's starboard settings",
			Usage:       "show",
			Permissions: crouter.PermLevelNone,
			Command:     commandStarboard,
		},
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "show",
		Aliases:     []string{"get"},
		Description: "Show the server's starboard settings",
		Usage:       "show",
		Permissions: crouter.PermLevelNone,
		Command:     commandStarboard,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "channel",
		Aliases:     []string{"ch"},
		Description: "Change the starboard channel",
		Usage:       "channel <channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardChannel,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "limit",
		Description: "Change the starboard limit",
		Usage:       "limit <int>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardLimit,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "emoji",
		Description: "Change the starboard emoji",
		Usage:       "emoji <emoji>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardEmoji,
	})

	router.AddCommand(&crouter.Command{
		Name:        "modlog",
		Description: "Set the modlog channel",
		Usage:       "modlog <channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandModLogChannel,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ban",
		Description: "Ban users by ID",
		Usage:       "ban <user ID> [reason]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Ban,
	})

	router.AddCommand(&crouter.Command{
		Name:        "mvc",
		Description: "Get the mvc majority",
		Usage:       "mvc [abstains]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.MVC,
	})

	gk := router.AddGroup(&crouter.Group{
		Name:        "gatekeeper",
		Aliases:     []string{"gb", "g"},
		Description: "Manage the server's gatekeeper",
		Command: &crouter.Command{
			Name:        "approve",
			Aliases:     []string{"a"},
			Description: "Approves a user in the gatekeeper",
			Usage:       "approve <user ID>",
			Permissions: crouter.PermLevelMod,
			Command:     commands.GkApprove,
		},
	})

	gk.AddCommand(&crouter.Command{
		Name:        "channel",
		Aliases:     []string{"chan", "ch"},
		Description: "Set the gatekeeper channel",
		Usage:       "channel <channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkChannel,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "message",
		Aliases:     []string{"msg"},
		Description: "Set the gatekeeper message",
		Usage:       "message <message>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkMessage,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "welcome-channel",
		Aliases:     []string{"welcome-ch", "wch"},
		Description: "Set the welcome channel",
		Usage:       "welcome-channel <channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandWelcomeChannel,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "welcome-message",
		Aliases:     []string{"welcome-msg", "wmsg"},
		Description: "Set the welcome message",
		Usage:       "welcome-message <message>",
		Permissions: crouter.PermLevelMod,
		Command:     commandWelcomeMessage,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "gk-roles",
		Description: "Set the gatekeeper roles",
		Usage:       "gk-roles <roles...>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkRoles,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "member-roles",
		Description: "Set the member roles",
		Usage:       "member-roles <roles...>",
		Permissions: crouter.PermLevelMod,
		Command:     commandMemberRoles,
	})
}

func addAdminCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "export",
		Description: "Export this server's mod logs",
		Usage:       "export",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandExport,
	})

	router.AddCommand(&crouter.Command{
		Name:        "prefix",
		Description: "Show this server's prefix or change it",
		Usage:       "prefix [new prefix]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandPrefix,
	})

	router.AddCommand(&crouter.Command{
		Name:        "modroles",
		Aliases:     []string{"mod-roles", "modrole"},
		Description: "List/modify this server's mod roles",
		Usage:       "modroles [add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandModRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "helperroles",
		Aliases:     []string{"helper-roles", "helper-role", "helperrole"},
		Description: "List/modify this server's helper roles",
		Usage:       "helperroles [add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandHelperRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "refreshmvc",
		Description: "Refresh the mvc role",
		Usage:       "refreshmvc",
		Permissions: crouter.PermLevelAdmin,
		Command:     commands.RefreshMVC,
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

	router.AddCommand(&crouter.Command{
		Name:        "restart",
		Aliases:     []string{"kill"},
		Description: "Stop the bot immediately (restarts with `sytemd`)",
		Usage:       "restart",
		Permissions: crouter.PermLevelOwner,
		Command:     commandKill,
	})

	router.AddCommand(&crouter.Command{
		Name:        "update",
		Description: "Update the bot in place",
		Usage:       "update",
		Permissions: crouter.PermLevelOwner,
		Command:     commandUpdate,
	})

	router.AddCommand(&crouter.Command{
		Name:        "error",
		Description: "Get an error by UUID",
		Usage:       "error <UUID>",
		Permissions: crouter.PermLevelOwner,
		Command:     commands.Error,
	})
}

func addAutoResponses() {
	router.AddResponse(&crouter.AutoResponse{
		Triggers: []string{"nyaa"},
		Response: func(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
			_, err = s.ChannelMessageSend(m.ChannelID, "*Nyaa nyaa~* <:meowpats:771890485978726431>")
			return err
		},
	})

	router.AddResponse(&crouter.AutoResponse{
		Triggers: []string{"covebot pronouns", "covebot, what are your pronouns?", "covebotn't pronouns", "covebotn't, what are your pronouns?"},
		Response: func(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
			_, err = s.ChannelMessageSend(m.ChannelID, "cat/cat/cats/cats/catself sound really good right now ~~actually, they *always* sound good~~")
			return err
		},
	})
}
