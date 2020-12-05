package main

import (
	"math/rand"
	"regexp"

	"github.com/Starshine113/covebotnt/commands/modcommands"
	"github.com/Starshine113/covebotnt/commands/modutilcommands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
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

func addAdminCommands() {
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
