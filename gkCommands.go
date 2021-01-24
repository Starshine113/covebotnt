package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func commandGkChannel(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<channel>",
			MissingArgs:  "<channel>",
		})
		return
	}
	channel, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	if channel.GuildID != ctx.Message.GuildID {
		ctx.CommandError(errors.New("channel must be in current guild"))
		return
	}
	err = ctx.Database.SetGatekeeperChannel(ctx.Message.GuildID, channel.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}
	_, err = ctx.Send("Set the gatekeeper channel to " + channel.Mention() + ".")
	return
}

func commandGkMessage(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<msg>",
			MissingArgs:  "<msg>",
		})
		return
	}
	msg := strings.Join(ctx.Args, " ")

	var msgB bytes.Buffer
	tmpl, err := template.New("gatekeeper").Parse(msg)
	if err != nil {
		sugar.Errorf("Error parsing template: %v", err)
		return err
	}
	exampleData := struct {
		User *discordgo.User
	}{
		User: ctx.Message.Author,
	}
	if err := tmpl.Execute(&msgB, exampleData); err != nil {
		sugar.Errorf("Error executing template: %v", err)
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Preview",
		Description: msgB.String(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Input",
				Value:  "```" + msg + "```",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Preview of the message that will be sent when a member joins:",
		Embed:   embed,
	})

	err = ctx.Database.SetGatekeeperMsg(ctx.Message.GuildID, msg)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	_, err = ctx.Send(fmt.Sprintf("%v Updated the gatekeeper welcome message.", crouter.SuccessEmoji))
	return
}

func commandWelcomeChannel(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<channel>",
			MissingArgs:  "<channel>",
		})
		return
	}
	channel, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	if channel.GuildID != ctx.Message.GuildID {
		ctx.CommandError(errors.New("channel must be in current guild"))
		return
	}
	err = ctx.Database.SetWelcomeChannel(ctx.Message.GuildID, channel.ID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}
	_, err = ctx.Send("Set the welcome channel to " + channel.Mention() + ".")
	return
}

func commandWelcomeMessage(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<msg>",
			MissingArgs:  "<msg>",
		})
		return
	}
	msg := strings.Join(ctx.Args, " ")

	var msgB bytes.Buffer
	tmpl, err := template.New("gatekeeper").Parse(msg)
	if err != nil {
		sugar.Errorf("Error parsing template: %v", err)
		return err
	}
	exampleData := struct {
		User *discordgo.User
	}{
		User: ctx.Message.Author,
	}
	if err := tmpl.Execute(&msgB, exampleData); err != nil {
		sugar.Errorf("Error executing template: %v", err)
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Preview",
		Description: msgB.String(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Input",
				Value:  "```" + msg + "```",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Preview of the message that will be sent when a member joins:",
		Embed:   embed,
	})

	err = ctx.Database.SetWelcomeMsg(ctx.Message.GuildID, msg)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = ctx.Database.RemoveFromGuildCache(ctx.Message.GuildID)
	if err != nil {
		return err
	}

	_, err = ctx.Send(fmt.Sprintf("%v Updated the welcome message.", crouter.SuccessEmoji))
	return
}
