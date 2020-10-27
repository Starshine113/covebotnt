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
		var perms error
		var embedDesc string
		var level int

		if perms = checkAdmin(ctx); perms == nil {
			level = 3
		} else if perms = checkModPerm(ctx, guildSettings); perms == nil {
			level = 2
		} else if perms = checkHelperPerm(ctx, guildSettings); perms == nil {
			level = 1
		} else {
			level = 0
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

		if level == 3 {
			embedDesc = adminCmdString + modCmdString + helperCmdString + userCmdString
		} else if level == 2 {
			embedDesc = modCmdString + helperCmdString + userCmdString
		} else if level == 1 {
			embedDesc = helperCmdString + userCmdString
		} else {
			embedDesc = userCmdString
		}

		_, err = ctx.Send(&discordgo.MessageEmbed{
			Title:       "Command list",
			Description: embedDesc,
			Color:       0x21a1a8,
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

	return
}

func cmdEmbed(cmd *Command) *discordgo.MessageEmbed {
	var aliases string
	var permLevel string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	switch cmd.Permissions {
	case PermLevelNone:
		permLevel = "None"
	case PermLevelHelper:
		permLevel = "Helper"
	case PermLevelMod:
		permLevel = "Moderator"
	case PermLevelAdmin:
		permLevel = "Admin"
	case PermLevelOwner:
		permLevel = "Owner"
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("`%v`", cmd.Name),
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
				Value:  permLevel,
				Inline: false,
			},
		},
	}

	return embed
}
