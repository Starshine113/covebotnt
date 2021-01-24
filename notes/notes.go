package notes

import "github.com/starshine-sys/covebotnt/crouter"

// Init adds all the commands from this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "Notes",
		Description: "Shows a user's notes",
		Usage:       "<user>",
		Permissions: crouter.PermLevelHelper,
		Command:     CommandNotes,
	})

	router.AddCommand(&crouter.Command{
		Name:        "SetNote",
		Aliases:     []string{"AddNote"},
		Description: "Set a note for a user",
		Usage:       "<user> <note>",
		Permissions: crouter.PermLevelHelper,
		Command:     CommandSetNote,
	})

	router.AddCommand(&crouter.Command{
		Name:        "DelNote",
		Aliases:     []string{"RemoveNote"},
		Description: "Remove a note by ID",
		Usage:       "<id>",
		Permissions: crouter.PermLevelMod,
		Command:     CommandDelNote,
	})
}
