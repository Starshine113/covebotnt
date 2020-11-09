package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// RoleInfo shows information about a role
func RoleInfo(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	role, err := ctx.ParseRole(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	members, err := ctx.Session.GuildMembers(ctx.Message.GuildID, "", 1000)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	created, err := discordgo.SnowflakeTimestamp(role.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var count int64 = 0
	for _, m := range members {
		for _, r := range m.Roles {
			if r == role.ID {
				count++
				break
			}
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Role info",
		Description: "`<@&" + role.ID + ">`",
		Color:       role.Color,
		Timestamp:   created.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Role created",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  role.ID,
				Inline: true,
			},
			{
				Name:   "Name",
				Value:  role.Name,
				Inline: true,
			},
			{
				Name:   "Mentionable",
				Value:  fmt.Sprint(role.Mentionable),
				Inline: true,
			},
			{
				Name:   "Color",
				Value:  fmt.Sprintf("#%v", strconv.FormatInt(int64(role.Color), 16)),
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprint(count),
				Inline: true,
			},
			{
				Name:   "Hoisted",
				Value:  fmt.Sprint(role.Hoist),
				Inline: true,
			},
			{
				Name:   "Position",
				Value:  fmt.Sprint(role.Position),
				Inline: true,
			},
			{
				Name:   "Created",
				Value:  fmt.Sprintf("%v ago", prettyDurationString(time.Since(created))),
				Inline: true,
			},
		},
	}

	_, err = ctx.Send(embed)
	return
}
