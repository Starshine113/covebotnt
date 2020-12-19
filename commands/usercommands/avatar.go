package usercommands

import (
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Avatar returns the user's, or the given member's, avatar
func Avatar(ctx *crouter.Ctx) (err error) {
	var u *discordgo.User
	if len(ctx.Args) >= 1 {
		u, err = ctx.ParseUser(strings.Join(ctx.Args, " "))
		if err != nil {
			ctx.CommandError(err)
			return nil
		}
	} else {
		u = ctx.Author
	}

	embed := &discordgo.MessageEmbed{
		Title: "Avatar for " + u.String(),
		Image: &discordgo.MessageEmbedImage{
			URL: u.AvatarURL("1024"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "User ID: " + u.ID,
		},
	}

	_, err = ctx.Send(embed)
	return
}
