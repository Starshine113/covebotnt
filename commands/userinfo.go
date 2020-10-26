package commands

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// UserInfo returns user info, formatted nicely
func UserInfo(ctx *cbctx.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	user, err := ctx.ParseMember(ctx.Author.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	if len(ctx.Args) == 1 {
		user, err = ctx.ParseMember(ctx.Args[0])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	msg, err := ctx.Send("Working, please wait...")
	if err != nil {
		return err
	}

	createdTime, err := discordgo.SnowflakeTimestamp(user.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	createdTime = createdTime.UTC()

	joinedTime, _ := user.JoinedAt.Parse()
	joinedTime = joinedTime.UTC()

	highestRoleName, highestRoleColour, perms, rls, err := getPerms(ctx.Session, ctx.Message.GuildID, user.User.ID)
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

	var roles string
	for _, role := range rls {
		roles += role.Mention() + ", "
	}
	roles = roles[:len(roles)-2]
	if len(roles) >= 1000 {
		roles = "Too many to list"
	}

	permString := strings.Join(permStrings, ", ")
	if len(permString) >= 950 {
		permString = permString[:950]
		permString += "..."
	}

	guildCreated, err := discordgo.SnowflakeTimestamp(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	guildCreated = guildCreated.UTC()
	timeSinceCreation := joinedTime.Sub(guildCreated)
	days, _ := math.Modf(timeSinceCreation.Hours() / 24)

	nickname := user.User.Username
	if user.Nick != "" {
		nickname = user.Nick
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
				Name:   "Avatar",
				Value:  fmt.Sprintf("[Link](%v)", user.User.AvatarURL("1024")),
				Inline: true,
			},
			{
				Name:   "Highest rank",
				Value:  highestRoleName,
				Inline: true,
			},
			{
				Name:   "Username",
				Value:  user.User.String(),
				Inline: true,
			},
			{
				Name:   "Nickname",
				Value:  nickname,
				Inline: true,
			},
			{
				Name:   "Created at",
				Value:  fmt.Sprintf("%v\n(%v ago)", createdTime.Format("Jan _2 2006, 15:04:05 MST"), prettyDurationString(time.Since(createdTime))),
				Inline: true,
			},
			{
				Name:   "Joined at",
				Value:  fmt.Sprintf("%v\n(%v ago)\n%v days after the server was created", joinedTime.Format("Jan _2 2006, 15:04:05 MST"), prettyDurationString(time.Since(joinedTime)), days),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("Roles (%v)", len(user.Roles)),
				Value:  roles,
				Inline: false,
			},
			{
				Name:   "Permissions",
				Value:  fmt.Sprintf("```%v```", permString),
				Inline: false,
			},
		},
	}

	content := ""
	_, err = ctx.Edit(msg, &discordgo.MessageEdit{
		Content: &content,
		Embed:   embed,
	})
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
	out += fmt.Sprintf("%v minutes", minutes)

	return
}

func getPerms(s *discordgo.Session, guildID, userID string) (highestRoleName string, highestRoleColour, perms int, roles sortRoleByPos, err error) {
	// get the member
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return
	}

	rls := make(sortRoleByPos, len(member.Roles))
	for i, x := range member.Roles {
		r, err := s.State.Role(guildID, x)
		if err != nil {
			return highestRoleName, highestRoleColour, perms, rls, err
		}
		rls[i] = r
	}
	sort.Sort(rls)

	for _, role := range rls {
		perms |= role.Permissions
	}
	highestRoleName = rls[0].Name
	for _, role := range rls {
		if role.Color != 0 {
			return highestRoleName, role.Color, perms, rls, nil
		}
	}
	return
}

type sortRoleByPos []*discordgo.Role

func (a sortRoleByPos) Len() int           { return len(a) }
func (a sortRoleByPos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortRoleByPos) Less(i, j int) bool { return a[i].Position > a[j].Position }
