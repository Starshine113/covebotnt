package triggers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
)

func cmdimport(ctx *crouter.Ctx) (err error) {
	var fileURL string

	if len(ctx.Message.Attachments) > 0 {
		fileURL = ctx.Message.Attachments[0].URL
	}

	if len(ctx.Args) > 0 {
		fileURL = ctx.Args[0]
	}

	if fileURL == "" {
		_, err = ctx.SendNoAddXHandler(crouter.ErrorEmoji + "No URL or attachment supplied.")
		return err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	triggers, err := ctx.Database.Triggers(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var e export

	if err = json.Unmarshal(b, &e); err != nil {
		ctx.CommandError(err)
		return nil
	}

	var added, modified int

	for _, t := range e.Triggers {
		if t == nil {
			continue
		}

		t.GuildID = ctx.Message.GuildID
		// check if the trigger replaces one with the same ID
		var id int
		for _, x := range triggers {
			if t.Snowflake == x.Snowflake {
				id = x.ID
			}
			if strings.ToLower(t.Trigger) == strings.ToLower(x.Trigger) {
				id = x.ID
			}
		}

		if id != 0 {
			_, err = ctx.Database.EditTrigger(ctx.Message.GuildID, id, t)
			if err != nil {
				_, err = ctx.CommandError(err)
				return err
			}
			modified++
		} else {
			_, err = ctx.Database.AddTrigger(t)
			if err != nil {
				_, err = ctx.CommandError(err)
				return err
			}
			added++
		}
	}

	if added == 0 && modified == 0 {
		_, err = ctx.SendNoAddXHandler(crouter.WarnEmoji + " No entries were added/updated.")
	} else {
		_, err = ctx.SendfNoAddXHandler("Added %v entries and updated %v entries.", added, modified)
	}
	return
}
