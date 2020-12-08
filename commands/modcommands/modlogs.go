package modcommands

import (
	"fmt"
	"time"

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

	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	var user *discordgo.User

	member, err := ctx.ParseMember(ctx.Args[0])
	if err == nil {
		user = member.User
	} else {
		user, err = ctx.Session.User(ctx.Args[0])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	logs, err := ctx.Database.GetModLogs(ctx.Message.GuildID, user.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	if len(logs) == 0 {
		_, err = ctx.SendfNoAddXHandler("**%v** has no log entries.", user.String())
		return err
	}

	logs = reverseLogs(logs)

	logEntries := make([][]*cbdb.ModLogEntry, 0)

	for i := 0; i < len(logs); i += 5 {
		end := i + 5

		if end > len(logs) {
			end = len(logs)
		}

		logEntries = append(logEntries, logs[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for i, s := range logEntries {
		x := make([]*discordgo.MessageEmbedField, 0)
		for _, t := range s {
			x = append(x, logField(ctx, t))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    user.String(),
				IconURL: user.AvatarURL("128"),
			},
			Title: "Mod logs",
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v | User ID: %v", i+1, len(logEntries), user.ID),
			},
			Fields:    x,
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     0x21a1a8,
		})
	}

	if len(embeds) == 1 {
		_, err = ctx.Send(embeds[0])
		return
	}
	_, err = ctx.PagedEmbed(embeds)
	return
}

func logField(ctx *crouter.Ctx, log *cbdb.ModLogEntry) (field *discordgo.MessageEmbedField) {
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
	}
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
