package admincommands

import (
	"fmt"

	"github.com/Starshine113/covebotnt/crouter"
)

// Init adds all the commands in this package to the router
func Init(router *crouter.Router) {
	router.AddCommand(&crouter.Command{
		Name:        "Prefix",
		Description: "Show this server's prefix or change it",
		Usage:       "[new prefix]",
		Permissions: crouter.PermLevelAdmin,
		Command:     Prefix,
	})

	router.AddCommand(&crouter.Command{
		Name:        "ModRoles",
		Aliases:     []string{"ModRole"},
		Description: "List/modify this server's mod roles",
		Usage:       "[add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     ModRoles,
	})

	router.AddCommand(&crouter.Command{
		Name:        "HelperRoles",
		Aliases:     []string{"HelperRole"},
		Description: "List/modify this server's helper roles",
		Usage:       "[add|remove <role>]",
		Permissions: crouter.PermLevelAdmin,
		Command:     HelperRoles,
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
