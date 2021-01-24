package modcommands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/starshine-sys/covebotnt/cbdb"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

type export struct {
	GuildID   string              `json:"guild_id"`
	Timestamp time.Time           `json:"timestamp"`
	InvokedBy string              `json:"invoked_by"`
	Notes     []*cbdb.Note        `json:"notes"`
	ModLogs   []*cbdb.ModLogEntry `json:"mod_logs"`
}

// Export ...
func Export(ctx *crouter.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	notes, err := ctx.Database.AllNotes(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	logs, err := ctx.Database.GetAllLogs(ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	exportStruct := export{
		GuildID:   ctx.Message.GuildID,
		Timestamp: time.Now(),
		InvokedBy: ctx.Author.ID,
		Notes:     notes,
		ModLogs:   logs,
	}

	exportB, _ := json.MarshalIndent(exportStruct, "", "  ")
	reader := bytes.NewReader(exportB)

	file := discordgo.File{
		Name:   fmt.Sprintf("export-%v-%v.json", ctx.Message.GuildID, time.Now().Format("2006-01-02-15_04_05")),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})

	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, crouter.SuccessEmoji)
	if err != nil {
		return nil
	}

	return
}
