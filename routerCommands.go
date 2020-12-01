package main

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/Starshine113/covebotnt/commands/modcommands"
	"github.com/Starshine113/covebotnt/commands/modutilcommands"
	"github.com/Starshine113/covebotnt/commands/ownercommands"
	"github.com/Starshine113/covebotnt/commands/usercommands"
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
		Command:     usercommands.Ping,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Invite",
		Description: "Send an invite link for the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Invite,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Enlarge",
		Aliases:     []string{"E", "Emote", "Emoji", "Enlorge"},
		Description: "Enlarges up to 10 custom emoji",
		Usage:       "<emoji...>",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Enlarge,
		Cooldown:    5 * time.Second,
	})

	router.AddCommand(&crouter.Command{
		Name:        "EmojiInfo",
		Aliases:     []string{"EI", "EmoteInfo"},
		Description: "Get detailed info about a custom emoji",
		Usage:       "<emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.EmojiInfo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Color",
		Aliases:     []string{"Colour"},
		Description: "Previews a color",
		Usage:       "<color>",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Color,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Avatar",
		Aliases:     []string{"Pfp", "A"},
		Description: "Show a user's avatar",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Avatar,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Snowflake",
		Aliases:     []string{"IDTime"},
		Description: "Get timestamps from the given ID(s)",
		Usage:       "[...IDs]",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Snowflake,
	})

	router.AddCommand(&crouter.Command{
		Name:        "UserInfo",
		Aliases:     []string{"I", "Info", "Whois", "Profile"},
		Description: "Show information about a user (or yourself)",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.UserInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "PKInfo",
		Aliases:     []string{"PKI", "PKUserInfo", "PKWhois", "PKProfile"},
		Description: "Show information about the user who sent a PluralKit-proxied message",
		Usage:       "<message ID>",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.PKUserInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RoleInfo",
		Aliases:     []string{"Ri"},
		Description: "Show information about a role",
		Usage:       "<role>",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.RoleInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ServerInfo",
		Aliases:     []string{"Si", "GuildInfo"},
		Description: "Show information about the current server",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.GuildInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "MajorityOf",
		Aliases:     []string{"Majority"},
		Description: "Get the majority of a number with abstains",
		Usage:       "<count> [abstains]",
		Permissions: crouter.PermLevelNone,
		Command:     modcommands.MajorityOf,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Hello",
		Aliases:     []string{"Hi", "Henlo", "Heya", "Heyo"},
		Description: "Say hi to the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     usercommands.Hello,
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
		Command:     modcommands.Warn,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModLogs",
		Description: "Show a user's modlogs",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     modcommands.ModLogs,
	})

	router.AddCommand(&crouter.Command{
		Name:        "LogMute",
		Description: "Log a user mute made with another bot",
		Usage:       "<user> [-d <duration>] [-hardmute] [reason]",
		Permissions: crouter.PermLevelHelper,
		Command:     modcommands.LogMute,
	})

	router.AddCommand(&crouter.Command{
		Name:        "BGC",
		Aliases:     []string{"BackgroundCheck"},
		Description: "Quickly check a user",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     modcommands.BGC,
	})
}

func addModCommands() {
	router.AddCommand(&crouter.Command{
		Name:        "Members",
		Description: "Show all members of a role",
		Usage:       "<role>",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.Members,
	})

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
		Command:     modutilcommands.Nickname,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Echo",
		Aliases:     []string{"Say", "Send"},
		Description: "Make the bot say something",
		Usage:       "[-ch <channel>] <message>",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.Echo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Steal",
		Aliases:     []string{"AddEmote", "AddEmoji"},
		Description: "Steal an emote by URL + name, or usage in message (with Nitro)",
		Usage:       "<emoji: url/emoji> [name]",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.Steal,
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
		Aliases:     []string{"React", "Reaction"},
		Description: "Change the starboard emoji",
		Usage:       "<emoji>",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.StarboardReact,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "ToggleSenderReact",
		Description: "Toggle whether or not a message's sender can react to it with the star emoji",
		Usage:       "",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.StarboardSenderCanReactToggle,
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
		Command:     modcommands.Ban,
	})

	router.AddCommand(&crouter.Command{
		Name:        "MVC",
		Description: "Get the mvc majority",
		Usage:       "[abstains]",
		Permissions: crouter.PermLevelMod,
		Command:     modcommands.MVC,
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
			Command:     modcommands.GkApprove,
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
		Name:        "Archive",
		Description: "Archive the current channel",
		Permissions: crouter.PermLevelAdmin,
		Command:     modutilcommands.Archive,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RefreshMVC",
		Description: "Refresh the mvc role",
		Usage:       "",
		Permissions: crouter.PermLevelAdmin,
		Command:     modutilcommands.RefreshMVC,
	})

	dm := router.AddGroup(&crouter.Group{
		Name:        "DM",
		Aliases:     []string{"DirectMessage"},
		Description: "DM a server user",
		Command: &crouter.Command{
			Name:        "Standard",
			Aliases:     []string{"Std"},
			Description: "Send a message to a user, showing the moderator's name",
			Usage:       "<user> <message>",
			Permissions: crouter.PermLevelAdmin,
			Command:     modutilcommands.DM,
		},
	})

	dm.AddCommand(&crouter.Command{
		Name:        "Anonymous",
		Aliases:     []string{"Anon"},
		Description: "Send a message to a user, hiding the moderator's name",
		Usage:       "<user> <message>",
		Permissions: crouter.PermLevelAdmin,
		Command:     modutilcommands.AnonDM,
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
		Name:        "AdminDM",
		Description: "Send any user sharing a server with the bot a message, including attachment",
		Usage:       "<user ID> <message>",
		Permissions: crouter.PermLevelOwner,
		Command:     ownercommands.AdminDM,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Guilds",
		Description: "Show all guilds",
		Usage:       "",
		Permissions: crouter.PermLevelOwner,
		Command:     ownercommands.Guilds,
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
		Regex: regexp.MustCompile("(?i)(?:(?:<@756804764112519188>|<@!756804764112519188>|covebot|covebotn't|covey)[ 's]+(?:(?:what are|what're|whatre|whatr|what is|whats|what's|what) (?:your|you're|youre|yours|youres|you'res|you) (?:pronoun)|pronoun)){1}"),
		Response: func(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
			n := rand.Intn(5)
			var msg string
			switch n {
			case 0:
				msg = "I'm really liking cat/cat/cats/cats/catself right now"
			case 1:
				msg = "cat/cat/cats/cats/catself sound good"
			case 2:
				msg = "I'll take cat/cat/cats/cats/catself"
			case 3:
				msg = "Let's go with cat/cat/cats/cats/catself"
			case 4:
				msg = "cat/cat/cats/cats/catself sound really good right now ~~actually, they *always* sound good~~"
			}
			_, err = s.ChannelMessageSend(m.ChannelID, msg)
			return err
		},
	})
}
