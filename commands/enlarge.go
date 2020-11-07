package commands

import (
	"fmt"
	"regexp"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Enlarge enlarges a custom emoji
func Enlarge(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) == 0 {
		_, err = ctx.CommandError(&crouter.ErrorNotEnoughArgs{NumRequiredArgs: 1, SuppliedArgs: 0})
		return err
	}

	emojiMatch, _ := regexp.Compile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")

	if emojiMatch.MatchString(ctx.Args[0]) {
		extension := ".png"
		groups := emojiMatch.FindStringSubmatch(ctx.Args[0])
		if groups[1] == "a" {
			extension = ".gif"
		}
		name := groups[2]
		url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)
		embed := &discordgo.MessageEmbed{
			Title:       ":" + name + ":",
			Description: fmt.Sprintf("<%v:%v:%v\\>", groups[1], groups[2], groups[3]),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "ID: " + groups[3],
			},
			Image: &discordgo.MessageEmbedImage{
				URL: url,
			},
		}
		_, err = ctx.Send(embed)
		return err
	}
	_, err = ctx.CommandError(&crouter.ErrorMissingRequiredArgs{RequiredArgs: "emoji", MissingArgs: "emoji"})
	return err
}
