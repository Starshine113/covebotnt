package usercommands

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Enlarge enlarges up to 10 custom emoji
func Enlarge(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		_, err = ctx.CommandError(&crouter.ErrorNotEnoughArgs{NumRequiredArgs: 1, SuppliedArgs: 0})
		return err
	}

	emojiMatch, _ := regexp.Compile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")

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
