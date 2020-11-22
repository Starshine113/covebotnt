package crouter

import (
	"fmt"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// NewRouter creates a Router object
func NewRouter(botOwners []string) *Router {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	router := &Router{
		BotOwners: botOwners,
		Cooldowns: cache,
	}

	router.AddCommand(&Command{
		Name:        "Help",
		Aliases:     []string{"Usage", "Hlep"},
		Description: "Show info about how to use the bot",
		Usage:       "[command]",
		Permissions: PermLevelNone,
		Command:     router.dummy,
	})

	return router
}

// dummy is used when a command isn't handled with the normal process
func (r *Router) dummy(ctx *Ctx) error {
	return nil
}

// AddCommand adds a command to the router
func (r *Router) AddCommand(cmd *Command) {
	cmd.Router = r
	if cmd.Cooldown == 0 {
		cmd.Cooldown = 500 * time.Millisecond
	}
	r.Commands = append(r.Commands, cmd)
}

// AddResponse adds an autoresponse to the router
func (r *Router) AddResponse(response *AutoResponse) {
	r.AutoResponses = append(r.AutoResponses, response)
}

// GetCommand gets a command by name
func (r *Router) GetCommand(name string) (c *Command) {
	for _, cmd := range r.Commands {
		if strings.ToLower(cmd.Name) == strings.ToLower(name) {
			return cmd
		}
		for _, a := range cmd.Aliases {
			if strings.ToLower(a) == strings.ToLower(name) {
				return cmd
			}
		}
	}
	return nil
}

// Respond checks if the message is any of the configured autoresponse triggers, and responds if it is
func (r *Router) Respond(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	for _, response := range r.AutoResponses {
		for _, trigger := range response.Triggers {
			if strings.ToLower(m.Content) == trigger {
				err = response.Response(s, m)
				return
			}
		}
		if response.Regex != nil {
			if response.Regex.MatchString(strings.ToLower(m.Content)) {
				err = response.Response(s, m)
				return
			}
		}
	}
	return
}

// Execute actually executes the router
func (r *Router) Execute(ctx *Ctx, guildSettings *structs.GuildSettings) (err error) {
	// add the guild settings to the additional parameters
	ctx.AdditionalParams["guildSettings"] = guildSettings
	ctx.GuildSettings = guildSettings

	help := r.GetCommand("help")
	if ctx.Match(append([]string{help.Name}, help.Aliases...)...) {
		err = r.Help(ctx, guildSettings)
		return
	}
	for _, g := range r.Groups {
		if ctx.Match(append([]string{g.Name}, g.Aliases...)...) || ctx.MatchRegexp(g.Regex) {
			if len(ctx.Args) == 0 {
				ctx.Command = ""
			} else {
				ctx.Command = ctx.Args[0]
			}
			if len(ctx.Args) > 1 {
				ctx.Args = ctx.Args[1:]
			} else {
				ctx.Args = []string{}
			}
			err = g.Execute(ctx, guildSettings)
			return
		}
	}
	for _, cmd := range r.Commands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) || ctx.MatchRegexp(cmd.Regex) {
			if len(ctx.Args) > 0 {
				if ctx.Args[0] == "help" || ctx.Args[0] == "usage" {
					ctx.Args[0] = ctx.Command
					err = r.Help(ctx, guildSettings)
					return
				}
			}
			ctx.Cmd = cmd
			if perms := ctx.Check(r.BotOwners); perms != nil {
				ctx.CommandError(perms)
				return nil
			}
			if cmd.Cooldown != time.Duration(0) {
				if _, e := r.Cooldowns.Get(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name)); e == nil {
					_, err = ctx.Sendf("The command `%v` can only be run once every **%v**.", cmd.Name, cmd.Cooldown.Round(time.Millisecond).String())
					return err
				}
				err = r.Cooldowns.SetWithTTL(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name), true, cmd.Cooldown)
				if err != nil {
					return err
				}
			}
			err = cmd.Command(ctx)
			return err
		}
	}
	return
}
