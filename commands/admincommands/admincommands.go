package admincommands

import (
	"fmt"

	"github.com/starshine-sys/covebotnt/crouter"
)

// Init adds all the commands in this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:            "ModRoles",
		Aliases:         []string{"ModRole"},
		Description:     "List/modify this server's mod roles",
		LongDescription: "Use `ModRoles Add` to add a moderator role.\nUse `ModRoles Remove` to remove a role.\nCalling `ModRoles` with no arguments will show the current list of moderator roles.",
		Usage:           "[add|remove <role>]",
		Permissions:     crouter.PermLevelAdmin,
		Command:         ModRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:            "HelperRoles",
		Aliases:         []string{"HelperRole"},
		Description:     "List/modify this server's helper roles",
		LongDescription: "Use `HelperRoles Add` to add a helper role.\nUse `HelperRoles Remove` to remove a role.\nCalling `HelperRoles` with no arguments will show the current list of helper roles.",
		Usage:           "[add|remove <role>]",
		Permissions:     crouter.PermLevelAdmin,
		Command:         HelperRoles,
	})

	p := router.AddGroup(&crouter.Group{
		Name:        "Prefix",
		Aliases:     []string{"Prefixes"},
		Description: "Manage the server's prefixes",
		Command: &crouter.Command{
			Name:        "Show",
			Description: "List the server's prefixes",
			Permissions: crouter.PermLevelNone,
			Command:     prefixShow,
			GuildOnly:   true,
		},
	})

	p.AddCommand(&crouter.Command{
		Name:        "Add",
		Description: "Add a prefix",
		Usage:       "<prefix>",
		Permissions: crouter.PermLevelAdmin,
		Command:     prefixAdd,
	})

	p.AddCommand(&crouter.Command{
		Name:        "Remove",
		Description: "Remove a prefix",
		Usage:       "<prefix>",
		Permissions: crouter.PermLevelAdmin,
		Command:     prefixRemove,
	})
}

func removeElement(s []string, i int) ([]string, error) {
	// s is [1,2,3,4,5,6], i is 2

	// perform bounds checking first to prevent a panic!
	if i >= len(s) || i < 0 {
		return nil, fmt.Errorf("Index is out of range. Index is %d with slice length %d", i, len(s))
	}

	// copy the last element (6) to index `i`. At this point,
	// `s` will be [1,2,6,4,5,6]
	s[i] = s[len(s)-1]
	// Remove the last element from the slice by truncating it
	// This way, `s` will now include all the elements from index 0
	// up to (but not including) the last element
	return s[:len(s)-1], nil
}
