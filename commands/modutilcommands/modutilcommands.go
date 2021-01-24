package modutilcommands

import "github.com/starshine-sys/covebotnt/crouter"

// Init adds all the commands from this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "Members",
		Description: "Show all members of a role",
		Usage:       "<role>",
		Permissions: crouter.PermLevelMod,
		Command:     Members,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Nickname",
		Aliases:     []string{"Nick"},
		Description: "Change the bot's nickname",
		Usage:       "<new nickname>",
		Permissions: crouter.PermLevelMod,
		Command:     Nickname,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Echo",
		Aliases:     []string{"Say", "Send"},
		Description: "Make the bot say something",
		Usage:       "[-ch <channel>] <message>",
		Permissions: crouter.PermLevelMod,
		Command:     Echo,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Steal",
		Aliases:     []string{"AddEmote", "AddEmoji"},
		Description: "Steal an emote by URL + name, or usage in message (with Nitro)",
		Usage:       "<emoji: url/emoji> [name]",
		Permissions: crouter.PermLevelMod,
		Command:     Steal,
	})

	router.AddCommand(&crouter.Command{
		Name:        "Archive",
		Description: "Archive the current channel",
		Permissions: crouter.PermLevelAdmin,
		Command:     Archive,
	})

	router.AddCommand(&crouter.Command{
		Name:        "RefreshMVC",
		Description: "Refresh the mvc role",
		Usage:       "",
		Permissions: crouter.PermLevelAdmin,
		Command:     RefreshMVC,
	})

	router.AddCommand(&crouter.Command{
		Name:        "CreateInvite",
		Description: "Create an invite for the given channel",
		Usage:       "[channel]",
		Permissions: crouter.PermLevelMod,
		Command:     invite,
	})

	dm(router)
}

func dm(router *crouter.Router) {
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
			Command:     DM,
		},
	})

	dm.AddCommand(&crouter.Command{
		Name:        "Anonymous",
		Aliases:     []string{"Anon"},
		Description: "Send a message to a user, hiding the moderator's name",
		Usage:       "<user> <message>",
		Permissions: crouter.PermLevelAdmin,
		Command:     AnonDM,
	})
}
