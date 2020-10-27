package crouter

import "github.com/Starshine113/covebotnt/cbctx"

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
