package triggers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func add(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var trigger, response string
	ctx.Args = strings.Split(strings.Join(ctx.Args, " "), "\n")
	if len(ctx.Args) < 2 {
		_, err = ctx.SendfNoAddXHandler("No response provided.")
		return err
	}

	trigger = ctx.Args[0]
	if len(ctx.Args) > 2 {
		response = strings.Join(ctx.Args[1:], "\n")
	} else {
		response = ctx.Args[1]
	}

	t, err := ctx.Database.AddTrigger(&cbdb.Trigger{
		GuildID:   ctx.Message.GuildID,
		Creator:   ctx.Author.ID,
		Trigger:   trigger,
		Response:  response,
		Snowflake: ctx.Bot.SnowflakeGen.Get(),
	})
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.Author.String(),
			IconURL: ctx.Author.AvatarURL("128"),
		},
		Title:       fmt.Sprintf("Trigger added (ID: `%v`)", t.Snowflake),
		Description: t.Response,
		Fields: []*discordgo.MessageEmbedField{{
			Name:   "Trigger",
			Value:  t.Trigger,
			Inline: false,
		}},
		Color:     0x21a1a8,
		Timestamp: t.Modified.UTC().Format(time.RFC3339),
	})
	return
}
