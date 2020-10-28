package crouter

import (
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
)

// PermLevel is the permission level of the command
type PermLevel int

const (
	// PermLevelNone can be used by all users
	PermLevelNone PermLevel = iota
	// PermLevelHelper requires a helper role
	PermLevelHelper
	// PermLevelMod requires a moderator role
	PermLevelMod
	// PermLevelAdmin requires admin perms
	PermLevelAdmin
	// PermLevelOwner requires the person to be the bot owner
	PermLevelOwner
)

// String gives the string representation of a permission level
func (p PermLevel) String() string {
	switch p {
	case PermLevelNone:
		return "[0] NONE"
	case PermLevelHelper:
		return "[1] HELPER"
	case PermLevelMod:
		return "[2] MODERATOR"
	case PermLevelAdmin:
		return "[3] ADMIN"
	case PermLevelOwner:
		return "[4] OWNER"
	}
	return fmt.Sprintf("PermLevel(%d)", p)
}

// Router is the command router
type Router struct {
	Commands []*Command
}

// Command is a single command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Permissions PermLevel
	Command     func(ctx *cbctx.Ctx) error
}
