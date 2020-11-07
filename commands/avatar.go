package commands

import (
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Avatar returns the user's, or the given member's, avatar
func Avatar(ctx *crouter.Ctx) (err error) {
	var user string
	if len(ctx.Args) >= 1 {
		user = strings.Join(ctx.Args, " ")
	} else {
		user = ctx.Author.ID
	}
	m, err := ctx.ParseMember(user)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	embed := &discordgo.MessageEmbed{
		Title: "Avatar for " + m.User.String(),
		Image: &discordgo.MessageEmbedImage{
			URL: m.User.AvatarURL("1024"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "User ID: " + m.User.ID,
		},
	}

	_, err = ctx.Send(embed)
	return
}
