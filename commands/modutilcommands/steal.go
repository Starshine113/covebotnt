package modutilcommands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Steal adds an emote by URL/ID
func Steal(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckArgRange(1, 2); err != nil {
		ctx.CommandError(err)
		return nil
	}

	var url string = ctx.Args[0]
	var name string

	emojiMatch, _ := regexp.Compile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
	if emojiMatch.MatchString(ctx.Args[0]) {
		extension := ".png"
		groups := emojiMatch.FindStringSubmatch(ctx.Args[0])
		if groups[1] == "a" {
			extension = ".gif"
		}
		name = groups[2]
		url = fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)
	}

	if len(ctx.Args) == 2 {
		name = ctx.Args[1]
	}

	resp, err := http.Get(url)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	defer resp.Body.Close()

	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	b64 := base64.StdEncoding.EncodeToString(image)

	if strings.HasSuffix(url, ".gif") {
		b64 = "data:image/gif;base64," + b64
	} else {
		b64 = "data:image/png;base64," + b64
	}

	emoji, err := ctx.Session.GuildEmojiCreate(ctx.Message.GuildID, name, b64, nil)
	if err != nil {
		switch err.(type) {
		case *discordgo.RESTError:
			_, err = ctx.Sendf("%v REST error occured. Common causes are: maximum emoji limit reached, emote above 256 KB size limit.", crouter.ErrorEmoji)
		default:
			_, err = ctx.CommandError(err)
		}
		return
	}
	_, err = ctx.Send(fmt.Sprintf("Added emoji %v with name %v.", emoji.MessageFormat(), emoji.Name))
	return err
}
