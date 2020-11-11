package main

import (
	"github.com/Starshine113/covebotnt/commands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/notes"
	"github.com/bwmarrin/discordgo"
)

func addUserCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "Ping",
		Description: "Ping pong!",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Ping,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Enlarge",
		Aliases:     []string{"E", "Emote", "Emoji", "Enlorge"},
		Description: "Enlarges a custom emoji",
		Usage:       "<emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Enlarge,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Color",
		Aliases:     []string{"Colour"},
		Description: "Previews a color",
		Usage:       "<color>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Color,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Avatar",
		Aliases:     []string{"Pfp", "A"},
		Description: "Show a user's avatar",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Avatar,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Snowflake",
		Aliases:     []string{"IDTime"},
		Description: "Get timestamps from the given ID(s)",
		Usage:       "[...IDs]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Snowflake,
	})

	router.AddCommand(&crouter.Command{
		Name:        "About",
		Description: "Show some info about the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     commands.About,
	})

	router.AddCommand(&crouter.Command{
		Name:        "UserInfo",
		Aliases:     []string{"I", "Info", "Whois", "Profile"},
		Description: "Show information about a user (or yourself)",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.UserInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RoleInfo",
		Aliases:     []string{"Ri"},
		Description: "Show information about a role",
		Usage:       "<role>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.RoleInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ServerInfo",
		Aliases:     []string{"Si", "GuildInfo"},
		Description: "Show information about the current server",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     commands.GuildInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "MajorityOf",
		Aliases:     []string{"Majority"},
		Description: "Get the majority of a number with abstains",
		Usage:       "<count> [abstains]",
		Permissions: crouter.PermLevelNone,
		Command:     commands.MajorityOf,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Hello",
		Aliases:     []string{"Hi", "Henlo", "Heya", "Heyo"},
		Description: "Say hi to the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Hello,
		GuildOnly:   true,
	})
}

func addHelperCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "Notes",
		Description: "Shows a user's notes",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     notes.CommandNotes,
	})

	router.AddCommand(&crouter.Command{
		Name:        "SetNote",
		Aliases:     []string{"AddNote"},
		Description: "Set a note for a user",
		Usage:       "<user> <note>",
		Permissions: crouter.PermLevelHelper,
		Command:     notes.CommandSetNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Warn",
		Description: "Warn a user",
		Usage:       "<user> <reason>",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.Warn,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModLogs",
		Description: "Show a user's modlogs",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.ModLogs,
	})

	router.AddCommand(&crouter.Command{
		Name:        "LogMute",
		Description: "Log a user mute made with another bot",
		Usage:       "<user> [-d <duration>] [-hardmute] [reason]",
		Permissions: crouter.PermLevelHelper,
		Command:     commands.LogMute,
	})
}

func addModCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "DelNote",
		Aliases:     []string{"RemoveNote"},
		Description: "Remove a note by ID",
		Usage:       "<id>",
		Permissions: crouter.PermLevelMod,
		Command:     notes.CommandDelNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Nickname",
		Aliases:     []string{"Nick"},
		Description: "Change the bot's nickname",
		Usage:       "<new nickname>",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Nickname,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Echo",
		Aliases:     []string{"Say", "Send"},
		Description: "Make the bot say something",
		Usage:       "[-ch <channel>] <message>",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Echo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Steal",
		Aliases:     []string{"AddEmote", "AddEmoji"},
		Description: "Steal an emote by URL + name, or usage in message (with Nitro)",
		Usage:       "<emoji: url/emoji> [name]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Steal,
	})

	starboard := router.AddGroup(&crouter.Group{
		Name:        "Starboard",
		Aliases:     []string{"Sb"},
		Description: "Manage the server's starboard",
		Command: &crouter.Command{
			Name:        "Show",
			Description: "Show the server's starboard settings",
			Usage:       "",
			Permissions: crouter.PermLevelNone,
			Command:     commandStarboard,
		},
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "Show",
		Aliases:     []string{"Get"},
		Description: "Show the server's starboard settings",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     commandStarboard,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "Channel",
		Aliases:     []string{"Ch"},
		Description: "Change the starboard channel",
		Usage:       "<channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardChannel,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "Limit",
		Description: "Change the starboard limit",
		Usage:       "<int>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardLimit,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "Emoji",
		Description: "Change the starboard emoji",
		Usage:       "<emoji>",
		Permissions: crouter.PermLevelMod,
		Command:     commandStarboardEmoji,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModLog",
		Description: "Set the modlog channel",
		Usage:       "<channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandModLogChannel,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Ban",
		Description: "Ban users by ID",
		Usage:       "<user ID> [reason]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.Ban,
	})

	router.AddCommand(&crouter.Command{
		Name:        "MVC",
		Description: "Get the mvc majority",
		Usage:       "[abstains]",
		Permissions: crouter.PermLevelMod,
		Command:     commands.MVC,
	})

	gk := router.AddGroup(&crouter.Group{
		Name:        "Gatekeeper",
		Aliases:     []string{"Gk", "G"},
		Description: "Manage the server's gatekeeper",
		Command: &crouter.Command{
			Name:        "Approve",
			Aliases:     []string{"a"},
			Description: "Approves a user in the gatekeeper",
			Usage:       "<user ID>",
			Permissions: crouter.PermLevelMod,
			Command:     commands.GkApprove,
		},
	})

	gk.AddCommand(&crouter.Command{
		Name:        "Channel",
		Aliases:     []string{"Chan", "Ch"},
		Description: "Set the gatekeeper channel",
		Usage:       "<channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkChannel,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "Message",
		Aliases:     []string{"Msg"},
		Description: "Set the gatekeeper message",
		Usage:       "<message>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkMessage,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "WelcomeChannel",
		Aliases:     []string{"WelcomeCh", "WCh"},
		Description: "Set the welcome channel",
		Usage:       "<channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandWelcomeChannel,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "WelcomeMessage",
		Aliases:     []string{"WelcomeMsg", "WMsg"},
		Description: "Set the welcome message",
		Usage:       "<message>",
		Permissions: crouter.PermLevelMod,
		Command:     commandWelcomeMessage,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "GkRoles",
		Description: "Set the gatekeeper roles",
		Usage:       "<roles...>",
		Permissions: crouter.PermLevelMod,
		Command:     commandGkRoles,
	})

	gk.AddCommand(&crouter.Command{
		Name:        "MemberRoles",
		Description: "Set the member roles",
		Usage:       "<roles...>",
		Permissions: crouter.PermLevelMod,
		Command:     commandMemberRoles,
	})
}

func addAdminCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "Export",
		Description: "Export this server's mod logs",
		Usage:       "",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandExport,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Prefix",
		Description: "Show this server's prefix or change it",
		Usage:       "[new prefix]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandPrefix,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModRoles",
		Aliases:     []string{"ModRole"},
		Description: "List/modify this server's mod roles",
		Usage:       "[add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandModRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "HelperRoles",
		Aliases:     []string{"HelperRole"},
		Description: "List/modify this server's helper roles",
		Usage:       "[add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     commandHelperRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RefreshMVC",
		Description: "Refresh the mvc role",
		Usage:       "",
		Permissions: crouter.PermLevelAdmin,
		Command:     commands.RefreshMVC,
	})
}

func addOwnerCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "SetStatus",
		Description: "Set the bot's status",
		Usage:       "<-replace/-append> [<status>|-clear]",
		Permissions: crouter.PermLevelOwner,
		Command:     commandSetStatus,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Restart",
		Aliases:     []string{"Kill"},
		Description: "Stop the bot immediately (restarts with `sytemd`)",
		Usage:       "",
		Permissions: crouter.PermLevelOwner,
		Command:     commandKill,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Update",
		Description: "Update the bot in place",
		Usage:       "",
		Permissions: crouter.PermLevelOwner,
		Command:     commandUpdate,
	})

	router.AddCommand(&crouter.Command{
		Name:        "AdminDM",
		Description: "Send any user sharing a server with the bot a message, including attachment",
		Usage:       "<user ID> <message>",
		Permissions: crouter.PermLevelOwner,
		Command:     commands.AdminDM,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Error",
		Description: "Get an error by UUID",
		Usage:       "<UUID>",
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
