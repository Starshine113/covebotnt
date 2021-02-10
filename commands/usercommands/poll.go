package usercommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/crouter"
)

var keycaps = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}

func poll(ctx *crouter.Ctx) (err error) {
	// have to strip whitespace but cant easily do that for a slice
	// so we get this shit
	args := strings.Split(strings.Join(strings.Fields(strings.Join(ctx.Args, " ")), " "), "|")

	if len(args) < 3 {
		_, err = ctx.Send(":x: Need at least a question and 2 options, split with `|`.")
		return err
	}

	question := args[0]
	options := args[1:]
	if len(options) > 10 {
		_, err = ctx.Send(":x: Too many options, maximum 10.")
		return err
	}

	var desc string
	for i, o := range options {
		desc += fmt.Sprintf("%v %v\n", keycaps[i], o)
	}

	if len(desc) > 2048 {
		_, err = ctx.Send(":x: Embed description too long.")
		return err
	}
	if len(question) > 256 {
		_, err = ctx.Send(":x: Question too long (maximum 256 characters)")
		return err
	}

	ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ctx.Message.ID)
	msg, err := ctx.Send(&discordgo.MessageEmbed{
		Title:       question,
		Description: desc,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("%s (%s)", ctx.Author, ctx.Author.ID),
			IconURL: ctx.Author.AvatarURL("128"),
		},
	})
	if err != nil {
		return err
	}

	for i := 0; i < len(options); i++ {
		err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, keycaps[i])
		if err != nil {
			return err
		}
	}
	return
}

func quickpoll(ctx *crouter.Ctx) (err error) {
	// indicate that were processing
	ctx.React("🔄")
	id := ctx.Message.ID

	// wait a second for pk
	time.Sleep(time.Second)

	m, err := pk.GetMessage(ctx.Message.ID)
	if err == nil {
		id = m.ID
	} else {
		ctx.Session.MessageReactionRemove(ctx.Channel.ID, ctx.Message.ID, "🔄", "@me")
	}

	err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, id, "👍")
	if err != nil {
		return err
	}
	err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, id, "👎")
	if err != nil {
		return err
	}
	return
}
