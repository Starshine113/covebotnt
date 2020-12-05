package usercommands

import (
	"time"

	"github.com/Starshine113/covebotnt/crouter"
)

// Init adds all the commands from this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "Ping",
		Description: "Ping pong!",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     Ping,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Invite",
		Description: "Send an invite link for the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     Invite,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Enlarge",
		Aliases:     []string{"E", "Emote", "Emoji", "Enlorge"},
		Description: "Enlarges up to 10 custom emoji",
		Usage:       "<emoji...>",
		Permissions: crouter.PermLevelNone,
		Command:     Enlarge,
		Cooldown:    5 * time.Second,
	})

	router.AddCommand(&crouter.Command{
		Name:        "EmojiInfo",
		Aliases:     []string{"EI", "EmoteInfo"},
		Description: "Get detailed info about a custom emoji",
		Usage:       "<emoji>",
		Permissions: crouter.PermLevelNone,
		Command:     EmojiInfo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Color",
		Aliases:     []string{"Colour"},
		Description: "Previews a color",
		Usage:       "<color>",
		Permissions: crouter.PermLevelNone,
		Command:     Color,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Avatar",
		Aliases:     []string{"Pfp", "A"},
		Description: "Show a user's avatar",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     Avatar,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Snowflake",
		Aliases:     []string{"IDTime"},
		Description: "Get timestamps from the given ID(s)",
		Usage:       "[...IDs]",
		Permissions: crouter.PermLevelNone,
		Command:     Snowflake,
	})

	router.AddCommand(&crouter.Command{
		Name:        "UserInfo",
		Aliases:     []string{"I", "Info", "Whois", "Profile"},
		Description: "Show information about a user (or yourself)",
		Usage:       "[user]",
		Permissions: crouter.PermLevelNone,
		Command:     UserInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "PKInfo",
		Aliases:     []string{"PKI", "PKUserInfo", "PKWhois", "PKProfile"},
		Description: "Show information about the user who sent a PluralKit-proxied message",
		Usage:       "<message ID>",
		Permissions: crouter.PermLevelNone,
		Command:     PKUserInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RoleInfo",
		Aliases:     []string{"Ri"},
		Description: "Show information about a role",
		Usage:       "<role>",
		Permissions: crouter.PermLevelNone,
		Command:     RoleInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ServerInfo",
		Aliases:     []string{"Si", "GuildInfo"},
		Description: "Show information about the current server",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     GuildInfo,
		GuildOnly:   true,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Hello",
		Aliases:     []string{"Hi", "Henlo", "Heya", "Heyo"},
		Description: "Say hi to the bot",
		Usage:       "",
		Permissions: crouter.PermLevelNone,
		Command:     Hello,
		GuildOnly:   true,
	})
}
