package commands

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// UserInfo returns user info, formatted nicely
func UserInfo(ctx *cbctx.Ctx) (err error) {
	user, err := ctx.ParseMember(ctx.Author.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	if len(ctx.Args) == 1 {
		user, err = ctx.ParseMember(ctx.Args[1])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}
	var roleList []string
	for _, role := range user.Roles {
		roleList = append(roleList, fmt.Sprintf("<@&%v>", role))
	}
	roles := strings.Join(roleList, ", ")
	if len(roles) >= 2000 {
		roles = "Too many to list"
	}

	createdTime, err := discordgo.SnowflakeTimestamp(user.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	createdTime = createdTime.UTC()

	joinedTime, _ := user.JoinedAt.Parse()
	joinedTime = joinedTime.UTC()

	highestRoleName, highestRoleColour, perms, err := getPerms(ctx.Session, ctx.Message.GuildID, user.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	var permStrings []string
	if perms&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		permStrings = append(permStrings, "Administrator")
	}
	if perms&discordgo.PermissionManageServer == discordgo.PermissionManageServer {
		permStrings = append(permStrings, "Manage Guild")
	}
	if perms&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels {
		permStrings = append(permStrings, "Manage Channels")
	}
	if perms&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles {
		permStrings = append(permStrings, "Manage Roles")
	}
	if perms&discordgo.PermissionBanMembers == discordgo.PermissionBanMembers {
		permStrings = append(permStrings, "Ban Members")
	}
	if perms&discordgo.PermissionMentionEveryone == discordgo.PermissionMentionEveryone {
		permStrings = append(permStrings, "Mention @everyone")
	}
	if perms&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages {
		permStrings = append(permStrings, "Manage Messages")
	}

	permString := strings.Join(permStrings, ", ")
	if len(permString) >= 950 {
		permString = permString[:950]
		permString += "..."
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.User.String(),
			IconURL: user.User.AvatarURL("256"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v | Created", user.User.ID),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.User.AvatarURL("512"),
		},
		Timestamp:   createdTime.Format(time.RFC3339),
		Color:       highestRoleColour,
		Description: fmt.Sprintf("User information for %v", user.Mention()),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Highest rank",
				Value:  highestRoleName,
				Inline: true,
			},
			{
				Name:   "Created at",
				Value:  fmt.Sprintf("%v\n(%v ago)", createdTime.Format("Jan _2 2006, 15:04:05 MST"), prettyDurationString(time.Since(createdTime))),
				Inline: true,
			},
			{
				Name:   "Joined at",
				Value:  fmt.Sprintf("%v\n(%v ago)", joinedTime.Format("Jan _2 2006, 15:04:05 MST"), prettyDurationString(time.Since(joinedTime))),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("Roles (%v)", len(user.Roles)),
				Value:  roles,
				Inline: false,
			},
			{
				Name:   "Permissions",
				Value:  fmt.Sprintf("`%v`", permString),
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)
	return
}

func prettyDurationString(duration time.Duration) (out string) {
	var days, hours, hoursFrac, minutes float64

	hours = duration.Hours()
	hours, hoursFrac = math.Modf(hours)
	minutes = hoursFrac * 60

	hoursFrac = math.Mod(hours, 24)
	days = (hours - hoursFrac) / 24
	hours = hours - (days * 24)
	minutes = minutes - math.Mod(minutes, 1)

	if days != 0 {
		out += fmt.Sprintf("%v days, ", days)
	}
	if hours != 0 {
		out += fmt.Sprintf("%v hours, ", hours)
	}
	out += fmt.Sprintf("%v minutes ago", minutes)

	return
}

func getPerms(s *discordgo.Session, guildID, userID string) (highestRoleName string, highestRoleColour, perms int, err error) {
	// get the guild
	guild, err := s.Guild(guildID)
	if err != nil {
		return
	}

	// get the member
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return
	}

	var highestRolePos int

	// iterate through all guild roles
	for _, r := range guild.Roles {
		// iterate through member roles
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				perms |= r.Permissions
				if r.Position > highestRolePos {
					highestRolePos = r.Position
					highestRoleName = r.Name
				}
			}
		}
		highestRolePos = 0
		// do it again
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				if r.Position > highestRolePos && r.Color != 0 {
					highestRolePos = r.Position
					highestRoleColour = r.Color
				}
			}
		}
	}
	return
}
