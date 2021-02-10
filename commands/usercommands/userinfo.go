package usercommands

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/starshine-sys/covebotnt/etc"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/starshine-sys/pkgo"
)

var pk = pkgo.NewSession(nil)

// PKUserInfo runs UserInfo with the user ID returned by
func PKUserInfo(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	msg, err := pk.GetMessage(ctx.Args[0])
	if err != nil {
		if err == pkgo.ErrMsgNotFound {
			_, err = ctx.Sendf("%v Message with ID `%v` not found by PK. Are you sure this message was proxied and hasn't been deleted?", crouter.ErrorEmoji, ctx.Args[0])
			return err
		}
		_, err = ctx.CommandError(err)
		return err
	}

	ctx.Args = []string{msg.Sender}
	err = UserInfo(ctx)
	return
}

func noMemberInfo(ctx *crouter.Ctx) (err error) {
	user, err := ctx.Session.User(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	ts, _ := discordgo.SnowflakeTimestamp(user.ID)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.String(),
			IconURL: user.AvatarURL("128"),
		},
		Description: user.ID,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("256"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v | Created", user.ID),
		},
		Timestamp: ts.UTC().Format(time.RFC3339),
		Color:     0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  user.String(),
				Inline: true,
			},
			{
				Name:   "Created",
				Value:  fmt.Sprintf("%v\n(%v ago)", ts.Format("Jan _2 2006, 15:04:05 MST"), etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(ts))),
				Inline: true,
			},
		},
	}

	_, err = ctx.Send(embed)

	return err
}

// UserInfo returns user info, formatted nicely
func UserInfo(ctx *crouter.Ctx) (err error) {
	user, err := ctx.ParseMember(ctx.Author.ID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if len(ctx.Args) == 1 {
		user, err = ctx.ParseMember(strings.Join(ctx.Args, " "))
		if err != nil {
			switch err.(type) {
			case *discordgo.RESTError:
				return noMemberInfo(ctx)
			default:
				_, err = ctx.CommandError(err)
				return err
			}
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
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == nil {
		if g.OwnerID == user.User.ID {
			permStrings = append(permStrings, "Server Owner")
			perms = discordgo.PermissionAll
		}
	}
	permStrings = append(permStrings, etc.PermStrings(perms)...)
	permString := strings.Join(permStrings, ", ")
	if len(permString) > 1000 {
		permString = permString[:1000] + "..."
	}
	if len(permString) == 0 {
		permString = "None"
	}

	var rolesSlice []string
	for _, role := range rls {
		rolesSlice = append(rolesSlice, role.Mention())
	}
	var roles string
	if len(rolesSlice) == 0 {
		roles = "No roles."
	} else {
		roles = strings.Join(rolesSlice, ", ")
	}
	if len(roles) >= 1000 {
		roles = "Too many to list"
	}

	botPerm := crouter.PermLevelNone

	for _, r := range rls {
		for _, helperRole := range ctx.GuildSettings.Moderation.HelperRoles {
			if r.ID == helperRole {
				botPerm = crouter.PermLevelHelper
				break
			}
		}
	}
	for _, r := range rls {
		for _, modRole := range ctx.GuildSettings.Moderation.ModRoles {
			if r.ID == modRole {
				botPerm = crouter.PermLevelHelper
				break
			}
		}
	}
	if perms&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		botPerm = crouter.PermLevelAdmin
	}
	for _, id := range ctx.Cmd.Router.BotOwners {
		if user.User.ID == id {
			botPerm = crouter.PermLevelOwner
		}
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
		Description: user.User.ID,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User information for",
				Value:  user.Mention(),
				Inline: false,
			},
			{
				Name:   "Avatar",
				Value:  fmt.Sprintf("[Link](%v)", user.User.AvatarURL("1024")),
				Inline: true,
			},
			{
				Name:   "Bot Permissions",
				Value:  "`" + botPerm.String() + "`",
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
				Value:  fmt.Sprintf("%v\n(%v ago)", createdTime.Format("Jan _2 2006, 15:04:05 MST"), etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(createdTime))),
				Inline: true,
			},
			{
				Name:   "Joined at",
				Value:  fmt.Sprintf("%v\n(%v ago)\n%v days after the server was created", joinedTime.Format("Jan _2 2006, 15:04:05 MST"), etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(joinedTime)), days),
				Inline: false,
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

func getPerms(s *discordgo.Session, guildID, userID string) (highestRoleName string, highestRoleColour int, perms int64, roles sortRoleByPos, err error) {
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
		if perms&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			perms |= discordgo.PermissionAll
			break
		}
		perms |= role.Permissions
	}
	highestRoleName = rls[0].Name
	for _, role := range rls {
		if role.Color != 0 {
			return highestRoleName, role.Color, perms, rls, nil
		}
	}
	return highestRoleName, 0, perms, rls, nil
}

type sortRoleByPos []*discordgo.Role

func (a sortRoleByPos) Len() int           { return len(a) }
func (a sortRoleByPos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortRoleByPos) Less(i, j int) bool { return a[i].Position > a[j].Position }
