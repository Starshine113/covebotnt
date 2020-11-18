package crouter

import (
	"regexp"
	"strings"

	"github.com/Starshine113/covebotnt/etc"
)

// MatchPrefix checks if the message matched any prefix
func (ctx *Ctx) MatchPrefix() bool {
	return etc.HasAnyPrefix(strings.ToLower(ctx.Message.Content), ctx.GuildPrefix, "<@"+ctx.BotUser.ID+">", "<@!"+ctx.BotUser.ID+">")
}

// Match checks if any of the given command aliases match
func (ctx *Ctx) Match(cmds ...string) bool {
	for _, cmd := range cmds {
		if strings.ToLower(ctx.Command) == strings.ToLower(cmd) {
			return true
		}
	}
	return false
}

// MatchRegexp checks if the command matches the given regex
func (ctx *Ctx) MatchRegexp(re *regexp.Regexp) bool {
	if re == nil {
		return false
	}
	return re.MatchString(strings.ToLower(ctx.Command))
}
