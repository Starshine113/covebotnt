package levels

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

// CommandLevel is the level/lvl command, used by end users
func CommandLevel(ctx *cbctx.Ctx) (err error) {
	if ctx.Message.GuildID == "" {
		return
	}

	user, err := ctx.ParseMember(ctx.Author.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	if len(ctx.Args) != 0 {
		user, err = ctx.ParseMember(ctx.Args[0])
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	}

	xp, err := ctx.BoltDb.GetXPForUser(user.User.ID, ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	lvl := math.Floor(math.Cbrt(float64(xp))) - 2
	nextLvl := lvl + 1
	if lvl < 0 {
		lvl = 0
		nextLvl = 3
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: user.User.AvatarURL("256"),
			Name:    user.User.String(),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.User.AvatarURL("512"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "User ID: " + user.User.ID,
		},
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       fmt.Sprintf("Level %v", lvl),
		Color:       0x4e24d8,
		Description: fmt.Sprintf("%v XP", xp),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("To next level (%v)", lvl+1),
				Value:  fmt.Sprintf("%v/%v XP", xp, math.Pow(nextLvl, 3)),
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)

	return
}

// CommandLeaderboard is the leaderboard command
func CommandLeaderboard(ctx *cbctx.Ctx) (err error) {
	if ctx.Message.GuildID == "" {
		return
	}

	var index int64 = 0
	if len(ctx.Args) > 0 {
		index, err = strconv.ParseInt(ctx.Args[0], 10, 0)
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
		index--
		if index < 0 {
			index = 0
		}
	}

	list, err := ctx.BoltDb.AllEntriesForGuild(ctx.Message.GuildID)
	if err != nil {
		return
	}

	indexMin := int(index * 20)
	indexMax := int((index - 1) * 20)
	if indexMax <= 0 {
		indexMax = 19
	}
	if len(list) <= indexMax {
		indexMax = len(list)
	}

	listSlice := list[indexMin:indexMax]

	var embedText string
	for i, v := range listSlice {
		lvl := math.Floor(math.Cbrt(float64(v.Xp))) - 2
		if lvl < 0 {
			lvl = 0
		}
		embedText += fmt.Sprintf("%v. <@%v> %v XP (level %v)\n", i+1, v.UserID, v.Xp, lvl)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Leaderboard",
		Color:       0x4e24d8,
		Description: embedText,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Showing entries %v-%v (out of %v)", indexMin+1, indexMax, len(list)),
		},
	}

	_, err = ctx.Send(embed)

	return
}
