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
		Aliases:     []string{"emote", "emoji"},
		Description: "Enlarges a custom emoji",
		Usage:       "enlarge <emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     commands.Enlarge,
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
