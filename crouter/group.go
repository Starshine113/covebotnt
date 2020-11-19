package crouter

import (
	"regexp"
	"strings"

	"github.com/Starshine113/covebotnt/structs"
)

// Group is a group of subcommands
type Group struct {
	Name        string
	Aliases     []string
	Regex       *regexp.Regexp
	Description string
	Command     *Command
	Subcommands []*Command
	Router      *Router
}

// AddGroup adds a group to the router
func (r *Router) AddGroup(group *Group) *Group {
	group.Router = r
	r.Groups = append(r.Groups, group)
	return r.GetGroup(group.Name)
}

// AddCommand adds a command to a group
func (g *Group) AddCommand(cmd *Command) {
	cmd.Router = g.Router
	g.Subcommands = append(g.Subcommands, cmd)
}

// GetGroup returns a group by name
func (r *Router) GetGroup(name string) (group *Group) {
	for _, group := range r.Groups {
		if group.Name == name {
			return group
		}
		for _, a := range group.Aliases {
			if a == name {
				return group
			}
		}
	}
	return nil
}

// GetCommand gets a command by name
func (g *Group) GetCommand(name string) (c *Command) {
	for _, cmd := range g.Subcommands {
		if strings.ToLower(cmd.Name) == strings.ToLower(name) {
			return cmd
		}
		for _, a := range cmd.Aliases {
			if strings.ToLower(a) == strings.ToLower(name) {
				return cmd
			}
		}
	}
	if strings.ToLower(g.Command.Name) == strings.ToLower(name) {
		return g.Command
	}
	for _, a := range g.Command.Aliases {
		if strings.ToLower(a) == strings.ToLower(name) {
			return g.Command
		}
	}
	return nil
}

// Execute executes any command given
func (g *Group) Execute(ctx *Ctx, guildSettings *structs.GuildSettings) (err error) {
	g.Subcommands = append(g.Subcommands, g.Command)
	for _, cmd := range g.Subcommands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) || ctx.MatchRegexp(cmd.Regex) {
			ctx.Cmd = cmd
			if perms := ctx.Check(g.Router.BotOwners); perms != nil {
				ctx.CommandError(perms)
				return nil
			}
			err = cmd.Command(ctx)
			return
		}
	}
	ctx.Cmd = g.Command
	if perms := ctx.Check(g.Router.BotOwners); perms != nil {
		ctx.CommandError(perms)
		return nil
	}
	err = g.Command.Command(ctx)
	return
}
