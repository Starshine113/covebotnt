package crouter

import (
	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
)

// Group is a group of subcommands
type Group struct {
	Name        string
	Aliases     []string
	Description string
	Command     *Command
	Subcommands []*Command
}

// AddGroup adds a group to the router
func (r *Router) AddGroup(group *Group) *Group {
	r.Groups = append(r.Groups, group)
	return r.GetGroup(group.Name)
}

// AddCommand adds a command to a group
func (g *Group) AddCommand(cmd *Command) {
	g.Subcommands = append(g.Subcommands, cmd)
}

// GetGroup returns a group by name
func (r *Router) GetGroup(name string) (group *Group) {
	for _, group := range r.Groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}

// Execute executes any command given
func (g *Group) Execute(ctx *cbctx.Ctx, guildSettings *structs.GuildSettings, ownerIDs []string) (err error) {
	for _, cmd := range g.Subcommands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) {
			ctx.TriggerTyping()
			if cmd.Permissions == PermLevelHelper {
				perms := checkHelperPerm(ctx, guildSettings)
				if perms != nil {
					ctx.CommandError(perms)
					return nil
				}
			} else if cmd.Permissions == PermLevelMod {
				perms := checkModPerm(ctx, guildSettings)
				if perms != nil {
					ctx.CommandError(perms)
					return nil
				}
			} else if cmd.Permissions == PermLevelAdmin {
				perms := checkAdmin(ctx)
				if perms != nil {
					ctx.CommandError(perms)
					return nil
				}
			} else if cmd.Permissions == PermLevelOwner {
				perms := checkOwner(ctx.Author.ID, ownerIDs)
				if perms != nil {
					ctx.CommandError(perms)
					return nil
				}
			}
			err = cmd.Command(ctx)
			return
		}
	}
	if g.Command.Permissions == PermLevelHelper {
		perms := checkHelperPerm(ctx, guildSettings)
		if perms != nil {
			ctx.CommandError(perms)
			return nil
		}
	} else if g.Command.Permissions == PermLevelMod {
		perms := checkModPerm(ctx, guildSettings)
		if perms != nil {
			ctx.CommandError(perms)
			return nil
		}
	} else if g.Command.Permissions == PermLevelAdmin {
		perms := checkAdmin(ctx)
		if perms != nil {
			ctx.CommandError(perms)
			return nil
		}
	} else if g.Command.Permissions == PermLevelOwner {
		perms := checkOwner(ctx.Author.ID, ownerIDs)
		if perms != nil {
			ctx.CommandError(perms)
			return nil
		}
	}
	err = g.Command.Command(ctx)
	return
}
