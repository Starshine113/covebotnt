package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/bwmarrin/discordgo"
)

type export struct {
	GuildID   string       `json:"guild_id"`
	Timestamp time.Time    `json:"timestamp"`
	InvokedBy string       `json:"invoked_by"`
	Notes     []*cbdb.Note `json:"notes"`
}

func commandExport(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	perms := checkModRole(ctx.Session, ctx.Author.ID, ctx.Message.GuildID, true)
	if perms != nil {
		commandError(perms, ctx.Session, ctx.Message)
		return nil
	}

	notes, err := ctx.Database.AllNotes(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	exportStruct := export{
		GuildID:   ctx.Message.GuildID,
		Timestamp: time.Now(),
		InvokedBy: ctx.Author.ID,
		Notes:     notes,
	}

	exportB, _ := json.MarshalIndent(exportStruct, "", "    ")
	reader := bytes.NewReader(exportB)

	file := discordgo.File{
		Name:   fmt.Sprintf("notes-export-%v-%v.json", ctx.Message.GuildID, time.Now().Format("2006-01-02-15_04_05")),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})

	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, cbctx.SuccessEmoji)
	if err != nil {
		return nil
	}

	return
}

type archive struct {
	GuildID   string               `json:"guild_id"`
	ChannelID string               `json:"channel_id"`
	Timestamp time.Time            `json:"timestamp"`
	Messages  []*discordgo.Message `json:"messages"`
}

func commandArchive(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	perms := checkModRole(ctx.Session, ctx.Author.ID, ctx.Message.GuildID, true)
	if perms != nil {
		commandError(perms, ctx.Session, ctx.Message)
		return nil
	}

	messages, err := ctx.Session.ChannelMessages(ctx.Message.ChannelID, 100, "", "", "")
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	exportB, _ := json.MarshalIndent(archive{
		GuildID:   ctx.Message.GuildID,
		ChannelID: ctx.Message.ChannelID,
		Timestamp: time.Now(),
		Messages:  messages,
	}, "", "    ")

	fmt.Println(string(exportB))
	return nil
}
