package modcommands

import (
	"fmt"
	"strconv"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"

	"github.com/bwmarrin/discordgo"
)

// ModLogs shows the moderation logs for the specified user
func ModLogs(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    0,
		})
		return nil
	}

	var page int64 = 0

	if len(ctx.Args) > 1 {
		page, err = strconv.ParseInt(ctx.Args[1], 0, 0)
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
		page = page - 1
		if page < 0 {
			page = 0
		}
	}

	msg, err := ctx.Send("Working, please wait...")
	if err != nil {
		return err
	}
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	logs, err := ctx.Database.GetModLogs(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	logs = reverseLogs(logs)

	logFields := make([]*discordgo.MessageEmbedField, 0)

	minRange := page * 10
	maxRange := ((page + 1) * 10)
	if int64(len(logs)) < maxRange {
		maxRange = int64(len(logs))
	}
	if int64(len(logs)) < minRange {
		minRange = int64(len(logs)) - 10
		if minRange < 0 {
			minRange = 0
		}
	}
	logSlice := logs[minRange:maxRange]

	for _, log := range logSlice {
		field, err := logField(ctx, log)
		if err != nil {
			return err
		}
		logFields = append(logFields, field)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.String(),
			IconURL: member.User.AvatarURL("256"),
		},
		Title: "Mod logs",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v-%v out of %v shown | ID: %v", minRange+1, minRange+int64(len(logSlice)), len(logs), member.User.ID),
		},
		Fields:    logFields,
		Timestamp: string(ctx.Message.Timestamp),
	}

	_, err = ctx.Send(embed)
	if err != nil {
		return err
	}
	err = ctx.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
	return
}

func logField(ctx *crouter.Ctx, log *cbdb.ModLogEntry) (field *discordgo.MessageEmbedField, err error) {
	mod, err := ctx.Bot.MemberCache.Get(log.GuildID, log.ModID)
	var logValue string
	if err == nil {
		logValue = fmt.Sprintf("**Responsible moderator:** %v\n**Reason:** %v", mod.User.String(), log.Reason)
	} else {
		logValue = fmt.Sprintf("**Responsible moderator:** %v\n**Reason:** %v", log.ModID, log.Reason)
	}

	return &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("#%v | %v | %v", log.ID, log.Type, log.Time.Format("2006-01-02")),
		Value:  logValue,
		Inline: false,
	}, nil
}

func reverseLogs(s []*cbdb.ModLogEntry) []*cbdb.ModLogEntry {
	a := make([]*cbdb.ModLogEntry, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}
