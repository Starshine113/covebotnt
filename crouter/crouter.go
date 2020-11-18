package crouter

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
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
	Commands      []*Command
	Groups        []*Group
	AutoResponses []*AutoResponse
}

// AutoResponse is a single autoresponse, intended for very simple responses to exact messages that don't match commands
type AutoResponse struct {
	Triggers []string
	Regex    *regexp.Regexp
	Response func(s *discordgo.Session, m *discordgo.MessageCreate) error
}

// Command is a single command
type Command struct {
	Name        string
	Aliases     []string
	Regex       *regexp.Regexp
	Description string
	Usage       string
	Permissions PermLevel
	Command     func(*Ctx) error
	GuildOnly   bool
}
