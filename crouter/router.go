package crouter

import (
	"strings"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// NewRouter creates a Router object
func NewRouter() *Router {
	router := &Router{}

	router.AddCommand(&Command{
		Name:        "help",
		Aliases:     []string{"usage", "hlep"},
		Description: "Show info about how to use the bot",
		Usage:       "help [command]",
		Permissions: PermLevelNone,
		Command:     router.dummy,
	})

	return router
}

// dummy is used when a command isn't handled with the normal process
func (r *Router) dummy(ctx *cbctx.Ctx) error {
	return nil
}

// AddCommand adds a command to the router
func (r *Router) AddCommand(cmd *Command) {
	r.Commands = append(r.Commands, cmd)
}

// AddResponse adds an autoresponse to the router
func (r *Router) AddResponse(response *AutoResponse) {
	r.AutoResponses = append(r.AutoResponses, response)
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
	}
	return
}

// Execute actually executes the router
func (r *Router) Execute(ctx *cbctx.Ctx, guildSettings *structs.GuildSettings, ownerIDs []string) (err error) {
	if ctx.Match("help") {
		ctx.TriggerTyping()
		err = r.Help(ctx, guildSettings)
		return
	}
	for _, cmd := range r.Commands {
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
			// add the guild settings to the additional parameters
			ctx.AdditionalParams["guildSettings"] = guildSettings

			err = cmd.Command(ctx)
			return err
		}
	}
	return
}
