package triggers

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
	Timestamp time.Time       `json:"timestamp"`
	Triggers  []*cbdb.Trigger `json:"triggers"`
}

func cmdexport(ctx *crouter.Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return err
	}

	triggers, err := ctx.Database.Triggers(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	exportStruct := export{
		Timestamp: time.Now(),
		Triggers:  triggers,
	}

	exportB, _ := json.MarshalIndent(exportStruct, "", "  ")
	reader := bytes.NewReader(exportB)

	file := discordgo.File{
		Name:   fmt.Sprintf("triggers-%v-%v.json", ctx.Message.GuildID, time.Now().Format("2006-01-02-15_04_05")),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})

	err = ctx.React(crouter.SuccessEmoji)
	if err != nil {
		return nil
	}

	return
}
