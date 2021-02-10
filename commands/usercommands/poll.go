package usercommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/crouter"

	flag "github.com/spf13/pflag"
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
	var reacts int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.IntVarP(&reacts, "options", "o", -1, "How many options to have")
	fs.Parse(ctx.Args)
	ctx.Args = fs.Args()

	// i cant be bothered to write actual error messages so well just do this
	if reacts < 2 {
		_, err = ctx.Send("less than 2 options? do you really think i can do something with that?")
		return err
	} else if reacts > 10 {
		_, err = ctx.Send("look buddy i can't help you with that, that's way too many choices, i can only count to 10 smh")
		return err
	}

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

	if reacts < 2 || reacts > 10 {
		err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, id, ":greentick:754647778390442025")
		if err != nil {
			return err
		}
		err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, id, ":redtick:754647803837415444")
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < reacts; i++ {
			err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, id, keycaps[i])
			if err != nil {
				return err
			}
		}
	}
	return
}
