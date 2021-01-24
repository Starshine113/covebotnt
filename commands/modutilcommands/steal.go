package modutilcommands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/starshine-sys/covebotnt/crouter"
)

var emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
var nameRegex = regexp.MustCompile("[\\w\\d]{2,32}")

// Steal adds an emote by URL/ID
func Steal(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckArgRange(1, 2); err != nil {
		ctx.CommandError(err)
		return nil
	}

	var url string = ctx.Args[0]
	var name string

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

	if !nameRegex.MatchString(name) {
		_, err = ctx.UsageEmbed("Name must be between 2 and 32 characters long, and only contain alphanumeric characters and underscores.", "")
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		_, err = ctx.UsageEmbed("Invalid URL (or you put the name and URL in the wrong order). A valid URL starts with `http://` or `https://`.", "")
		return
	}

	if !strings.Contains(url, ".gif") && !strings.Contains(url, ".png") && !strings.Contains(url, ".jpg") && !strings.Contains(url, ".jpeg") {
		_, err = ctx.UsageEmbed("The given URL doesn't point to an image (should end with .png or .gif).", "")
		return
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
		_, err = ctx.CommandError(err)
		return
	}
	_, err = ctx.Send(fmt.Sprintf("Added emoji %v with name %v.", emoji.MessageFormat(), emoji.Name))
	return err
}
