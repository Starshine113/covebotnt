package crouter

import (
	"fmt"
	"strings"

	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// Help is the help command
func (r *Router) Help(ctx *Ctx, guildSettings *structs.GuildSettings) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	if len(ctx.Args) == 0 {
		level := 0
		permLevel := PermLevelNone

		if err = checkOwner(ctx.Author.ID, r.BotOwners); err == nil {
			level = 4
			permLevel = PermLevelOwner
		} else if err = checkAdmin(ctx); err == nil {
			level = 3
			permLevel = PermLevelAdmin
		} else if err = checkModPerm(ctx, guildSettings); err == nil {
			level = 2
			permLevel = PermLevelMod
		} else if err = checkHelperPerm(ctx, guildSettings); err == nil {
			level = 1
			permLevel = PermLevelHelper
		}

		var ownerCmdString, adminCmdString, modCmdString, helperCmdString, userCmdString string
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
			case PermLevelOwner:
				ownerCmdString += fmt.Sprintf("`%v`: %v\n", cmd.Name, cmd.Description)
			}
		}
		var groupCmds string
		for _, g := range r.Groups {
			groupCmds += fmt.Sprintf("`%v`: %v\n", g.Name, g.Description)
		}

		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       ctx.BotUser.Username + " help",
			Description: "Use ⬅️ ➡️ to navigate between pages, and ❌ to delete this message.",
			Color:       0x21a1a8,
			Fields: []*discordgo.MessageEmbedField{{
				Name:   "** **",
				Value:  fmt.Sprintf("This server's prefix is `%v`.\nYou can also mention the bot (%v) to invoke commands.\nYour bot permission level is `%v`.", ctx.GuildPrefix, ctx.BotUser.Mention(), permLevel.String()),
				Inline: false,
			}},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Use `help <cmd>` for more information on a command",
			},
		}, &discordgo.MessageEmbed{
			Title:       "Groups",
			Description: groupCmds,
			Color:       0x21a1a8,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Use `help <cmd>` for more information on a command",
			},
		}, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("User commands (%v)", len(strings.Split(userCmdString, "\n"))-1),
			Description: userCmdString,
			Color:       0x21a1a8,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Use `help <cmd>` for more information on a command",
			},
		})

		if level >= 1 {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Helper commands (%v)", len(strings.Split(helperCmdString, "\n"))-1),
				Description: helperCmdString,
				Color:       0x21a1a8,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Use `help <cmd>` for more information on a command",
				},
			})
		}
		if level >= 2 {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Mod commands (%v)", len(strings.Split(modCmdString, "\n"))-1),
				Description: modCmdString,
				Color:       0x21a1a8,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Use `help <cmd>` for more information on a command",
				},
			})
		}
		if level >= 3 {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Admin commands (%v)", len(strings.Split(adminCmdString, "\n"))-1),
				Description: adminCmdString,
				Color:       0x21a1a8,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Use `help <cmd>` for more information on a command",
				},
			})
		}

		_, err = ctx.PagedEmbed(embeds)
		return err
	}

	var cmd *Command
	g := r.GetGroup(ctx.Args[0])
	if g != nil {
		if len(ctx.Args) == 1 {
			_, err = ctx.Send(ctx.GroupEmbed(g))
			return
		}
		cmd = g.GetCommand(ctx.Args[1])
		if cmd != nil {
			_, err = ctx.Send(ctx.GroupCmdEmbed(g, cmd))
			return
		}
	}
	cmd = r.GetCommand(ctx.Args[0])
	if cmd != nil {
		_, err = ctx.Send(ctx.CmdEmbed(cmd))
		return
	}

	_, err = ctx.Send(fmt.Sprintf("%v Invalid command or group provided:\n> `%v` is not a known command, group or alias.", ErrorEmoji, ctx.Args[0]))

	return
}

// GroupEmbed ...
func (ctx *Ctx) GroupEmbed(g *Group) *discordgo.MessageEmbed {
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
				Value:  fmt.Sprintf("```%v%v %v```", ctx.GuildPrefix, strings.ToLower(g.Command.Name), g.Command.Usage),
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

// GroupCmdEmbed ...
func (ctx *Ctx) GroupCmdEmbed(g *Group, cmd *Command) *discordgo.MessageEmbed {
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
				Value:  fmt.Sprintf("```%v%v %v %v```", ctx.GuildPrefix, strings.ToLower(g.Name), strings.ToLower(cmd.Name), cmd.Usage),
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

// CmdEmbed ...
func (ctx *Ctx) CmdEmbed(cmd *Command) *discordgo.MessageEmbed {
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
				Value:  fmt.Sprintf("```%v%v %v```", ctx.GuildPrefix, strings.ToLower(cmd.Name), cmd.Usage),
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
