package commands

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Starshine113/covebotnt/crouter"
)

// Color shows a preview of the given color code
func Color(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	ctx.Args[0] = strings.TrimPrefix(ctx.Args[0], "#")
	clr, err := strconv.ParseInt(ctx.Args[0], 16, 0)
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Color preview",
		Description: "#" + ctx.Args[0],
		Color:       int(clr),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://fakeimg.pl/256x256/" + ctx.Args[0] + "/?text=%20",
		},
	}

	_, err = ctx.Send(embed)
	return
}
