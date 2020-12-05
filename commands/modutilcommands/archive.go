package modutilcommands

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/flagparser"
	"github.com/bwmarrin/discordgo"
)

type arc struct {
	Version   int                `json:"version"`
	InvokedBy *discordgo.User    `json:"invoked_by"`
	Channel   *discordgo.Channel `json:"channel"`
	Messages  []*msg             `json:"messages"`
}

type msg struct {
	ID          string                         `json:"id"`
	Content     string                         `json:"content"`
	Timestamp   discordgo.Timestamp            `json:"timestamp"`
	Author      author                         `json:"author"`
	Attachments []*discordgo.MessageAttachment `json:"attachments"`
	Embeds      []*discordgo.MessageEmbed      `json:"embeds"`
	Reactions   []*discordgo.MessageReactions  `json:"reactions"`
}

type author struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	AvatarURL     string `json:"avatar_url"`
}

// Archive ...
func Archive(ctx *crouter.Ctx) (err error) {
	archive := arc{InvokedBy: ctx.Author, Channel: ctx.Channel}

	fp, _ := flagparser.NewFlagParser(flagparser.Bool("gzip", "gz"), flagparser.String("out", "o", "output"))

	args, err := fp.Parse(ctx.Args)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	var gz bool
	if args["gzip"].(bool) {
		gz = true
	}
	out := ctx.Channel.ID
	if args["out"].(string) != "" {
		channel, err := ctx.ParseChannel(args["out"].(string))
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
		out = channel.ID
	}

	ctx.Send("Working, please wait...")
	if err = ctx.Session.ChannelTyping(out); err != nil {
		return err
	}

	messages := make([]*msg, 0)
	var before string
	for {
		msgs, err := ctx.Session.ChannelMessages(ctx.Channel.ID, 100, before, "", "")
		fmt.Printf("Messages before ID %v, got %v messages\n", before, len(msgs))
		if len(msgs) == 0 {
			break
		}
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
		before = msgs[len(msgs)-1].ID
		for _, m := range msgs {
			messages = append(messages, &msg{
				ID:        m.ID,
				Content:   m.Content,
				Timestamp: m.Timestamp,
				Author: author{
					ID:            m.Author.ID,
					Username:      m.Author.Username,
					Discriminator: m.Author.Discriminator,
					Bot:           m.Author.Bot,
					AvatarURL:     m.Author.AvatarURL(""),
				},
				Attachments: m.Attachments,
				Embeds:      m.Embeds,
				Reactions:   m.Reactions,
			})
		}
	}

	archive.Messages = messages

	b, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	fn := fmt.Sprintf("export-%v-%v.json", ctx.Channel.Name, ctx.Message.Timestamp)

	var buf *bytes.Buffer
	if gz {
		buf = new(bytes.Buffer)
		zw := gzip.NewWriter(buf)
		zw.Name = fn
		_, err = zw.Write(b)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
		err = zw.Close()
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
		fn = fn + ".gz"
	} else {
		buf = bytes.NewBuffer(b)
	}

	file := discordgo.File{
		Name:   fn,
		Reader: buf,
	}

	_, err = ctx.Session.ChannelMessageSendComplex(out, &discordgo.MessageSend{
		Content: fmt.Sprintf("Done!\n> Archive of %v (#%v), invoked by %v at %v.", ctx.Channel.Mention(), ctx.Channel.Name, ctx.Author.String(), ctx.Message.Timestamp),
		Files:   []*discordgo.File{&file},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	})
	return err
}
