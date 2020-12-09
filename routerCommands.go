package main

import (
	"github.com/Starshine113/covebotnt/commands/modcommands"
	"github.com/Starshine113/covebotnt/commands/modutilcommands"
	"github.com/Starshine113/covebotnt/crouter"
)

func addStarboardCommands() {
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

	starboard.AddCommand(&crouter.Command{
		Name:        "Blacklist",
		Description: "Show the current starboard blacklist",
		Usage:       "",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.SbBlacklist,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "BlacklistAdd",
		Description: "Add a channel to the blacklist",
		Usage:       "",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.SbBlacklistAdd,
	})

	starboard.AddCommand(&crouter.Command{
		Name:        "BlacklistRemove",
		Aliases:     []string{"BlacklistDelete"},
		Description: "Remove a channel from the blacklist",
		Usage:       "",
		Permissions: crouter.PermLevelMod,
		Command:     modutilcommands.SbBlacklistRemove,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModLog",
		Description: "Set the modlog channel",
		Usage:       "<channel>",
		Permissions: crouter.PermLevelMod,
		Command:     commandModLogChannel,
	})
}

func addGkCommands() {
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
}
