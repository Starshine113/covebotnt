package usercommands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

type cmdList []*crouter.Command

func (c cmdList) Len() int      { return len(c) }
func (c cmdList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cmdList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

// Details ...
func Details(ctx *crouter.Ctx) (err error) {
	var cmds cmdList
	cmds = append(cmds, ctx.Cmd.Router.Commands...)
	sort.Sort(cmds)
	cmdSlices := make([][]*crouter.Command, 0)

	for i := 0; i < len(cmds); i += 5 {
		end := i + 5

		if end > len(cmds) {
			end = len(cmds)
		}

		cmdSlices = append(cmdSlices, cmds[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)
	for _, c := range cmdSlices {
		embeds = append(embeds, detailEmbed(c))
	}

	msg, err := ctx.PagedEmbed(embeds)

	ctx.AdditionalParams["cmdList"] = cmdSlices

	for i, e := range emoji {
		if err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, e); err != nil {
			return
		}

		index := i
		ctx.AddReactionHandlerOnce(msg.ID, e, func(ctx *crouter.Ctx) {
			ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
			cmdList := ctx.AdditionalParams["cmdList"].([][]*crouter.Command)

			cmdSlice := cmdList[ctx.AdditionalParams["page"].(int)]
			ctx.Send(ctx.CmdEmbed(cmdSlice[index]))
		})
	}

	return
}

var emoji = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

func detailEmbed(cmds []*crouter.Command) *discordgo.MessageEmbed {
	if len(cmds) > 5 {
		return nil
	}

	commands := make([]string, 0)

	for i, cmd := range cmds {
		commands = append(commands, fmt.Sprintf("%v `%v`: %v", emoji[i], cmd.Name, cmd.Description))
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Details",
		Description: strings.Join(commands, "\n\n"),
		Fields: []*discordgo.MessageEmbedField{{
			Name:   "Usage",
			Value:  "Use the arrows to switch pages.",
			Inline: false,
		}},
		Color: 0x21a1a8,
	}
	return embed
}
