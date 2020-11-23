package usercommands

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

var emojiMatch *regexp.Regexp

// Enlarge enlarges up to 10 custom emoji
func Enlarge(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		_, err = ctx.CommandError(&crouter.ErrorNotEnoughArgs{NumRequiredArgs: 1, SuppliedArgs: 0})
		return err
	}

	if emojiMatch == nil {
		emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
	}

	var attachments []*discordgo.File
	if len(ctx.Args) > 5 {
		ctx.Args = ctx.Args[:5]
	}
	for _, a := range ctx.Args {
		if emojiMatch.MatchString(a) {
			extension := ".png"
			groups := emojiMatch.FindStringSubmatch(a)
			if groups[1] == "a" {
				extension = ".gif"
			}
			name := groups[2]
			url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)
			resp, err := http.Get(url)
			if err != nil {
				_, err = ctx.CommandError(err)
				return err
			}
			defer resp.Body.Close()

			file := discordgo.File{
				Name:   name + extension,
				Reader: resp.Body,
			}

			attachments = append(attachments, &file)
		}
	}

	_, err = ctx.Send(&discordgo.MessageSend{Files: attachments})
	return
}

// EmojiInfo ...
func EmojiInfo(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	if emojiMatch == nil {
		emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
	}

	e := ctx.Args[0]

	if !emojiMatch.MatchString(e) {
		_, err = ctx.CommandError(&crouter.ErrorMissingRequiredArgs{RequiredArgs: "emoji", MissingArgs: "emoji"})
		return err
	}

	extension := ".png"
	groups := emojiMatch.FindStringSubmatch(e)
	if groups[1] == "a" {
		extension = ".gif"
	}
	name := groups[2]
	url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)

	created, err := discordgo.SnowflakeTimestamp(groups[3])
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       ":" + name + "\\:",
		Description: fmt.Sprintf("<%v:%v:%v\\>", groups[1], groups[2], groups[3]),
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
		Color: 0x21a1a8,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v | Emoji created", groups[3]),
		},
		Timestamp: created.Format(time.RFC3339),
	}

	_, err = ctx.Send(embed)
	return
}
