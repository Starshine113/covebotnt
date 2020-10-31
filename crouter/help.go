package crouter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// Help is the help command
func (r *Router) Help(ctx *cbctx.Ctx, guildSettings *structs.GuildSettings) (err error) {
	if len(ctx.Args) == 0 {
		level := 0

		if err = checkAdmin(ctx); err == nil {
			level = 3
		} else if err = checkModPerm(ctx, guildSettings); err == nil {
			level = 2
		} else if err = checkHelperPerm(ctx, guildSettings); err == nil {
			level = 1
		}

		var adminCmdString, modCmdString, helperCmdString, userCmdString string
		for _, cmd := range r.Commands {
			switch cmd.Permissions {
			case PermLevelNone:
				userCmdString += fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Description)
			case PermLevelHelper:
				helperCmdString += fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Description)
			case PermLevelMod:
				modCmdString += fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Description)
			case PermLevelAdmin:
				adminCmdString += fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Description)
			}
		}
		var groupCmds string
		for _, g := range r.Groups {
			groupCmds += fmt.Sprintf("`%v`: %v\n", g.Name, g.Description)
		}

		embedDesc := userCmdString
		if level == 1 {
			embedDesc += helperCmdString
		} else if level == 2 {
			embedDesc += modCmdString
		}
		if level == 3 {
			embedDesc += adminCmdString
		}

		_, err = ctx.Send(&discordgo.MessageEmbed{
			Title:       "Command list",
			Description: embedDesc,
			Color:       0x21a1a8,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Groups",
					Value:  groupCmds,
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Use `help <cmd>` for more information on a command",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	for _, cmd := range r.Commands {
		aliases := []string{cmd.Name}
		aliases = append(aliases, cmd.Aliases...)
		for _, alias := range aliases {
			if ctx.Args[0] == alias {
				_, err = ctx.Send(cmdEmbed(cmd))
				return err
			}
		}
	}
	for _, g := range r.Groups {
		aliases := []string{g.Name}
		aliases = append(aliases, g.Aliases...)
		for _, alias := range aliases {
			if ctx.Args[0] == alias {
				if len(ctx.Args) == 1 {
					_, err = ctx.Send(groupEmbed(g))
					return
				}
				for _, cmd := range g.Subcommands {
					for _, alias := range append([]string{cmd.Name}, cmd.Aliases...) {
						if ctx.Args[1] == alias {
							_, err = ctx.Send(groupCmdEmbed(g, cmd))
							return
						}
					}
				}
			}
		}
	}

	_, err = ctx.Send(fmt.Sprintf("%v Invalid command or group provided:\n> `%v` is not a known command, group or alias.", cbctx.ErrorEmoji, ctx.Args[0]))

	return
}

func groupEmbed(g *Group) *discordgo.MessageEmbed {
	var aliases string
	if g.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(g.Aliases, ", ")
	}

	var subCmds []string
	for _, cmd := range g.Subcommands {
		subCmds = append(subCmds, cmd.Name)
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(g.Name)),
		Description: g.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Subcommands",
				Value:  fmt.Sprintf("```%v```", strings.Join(subCmds, "\n")),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```\n** **", aliases),
				Inline: false,
			},
			{
				Name:   "Default command",
				Value:  g.Command.Description,
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v```", g.Command.Usage),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + g.Command.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

func groupCmdEmbed(g *Group, cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v %v```", strings.ToUpper(g.Name), strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v %v```", g.Name, cmd.Usage),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```", aliases),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + cmd.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}

func cmdEmbed(cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v```", cmd.Usage),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```", aliases),
				Inline: false,
			},
			{
				Name:   "Permission level",
				Value:  "```" + cmd.Permissions.String() + "```",
				Inline: false,
			},
		},
	}

	return embed
}
