package usercommands

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/starshine-sys/covebotnt/etc"
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

	var count int = 0
	for _, m := range members {
		for _, r := range m.Roles {
			if r == role.ID {
				count++
				break
			}
		}
	}

	permString := strings.Join(etc.PermStrings(role.Permissions), ", ")
	if len(permString) > 1000 {
		permString = permString[:1000] + "..."
	}
	if len(permString) == 0 {
		permString = "None"
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
				Name:   "Position",
				Value:  fmt.Sprint(role.Position),
				Inline: true,
			},
			{
				Name:   "Other info",
				Value:  fmt.Sprintf("Mentionable: %v\nHoisted: %v", role.Mentionable, role.Hoist),
				Inline: true,
			},
			{
				Name: "Percentage of users",
				Value: fmt.Sprintf("%v/%v (%v%%)", count, len(members),
					math.Round(float64(count)/float64(len(members))*100)),
				Inline: true,
			},
			{
				Name:   "Created",
				Value:  fmt.Sprintf("%v ago", etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(created))),
				Inline: true,
			},
			{
				Name:   "Permissions",
				Value:  "```" + permString + "```",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)
	return
}

func trinaryOperationThing(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}
