package triggers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
)

func remove(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	triggers, err := ctx.Database.Triggers(ctx.Message.GuildID)

	if len(triggers) == 0 {
		_, err = ctx.SendNoAddXHandler("There are no registered triggers.")
		return
	}

	var id int

	// try parsing an int
	if i, err := strconv.Atoi(ctx.Args[0]); err == nil {
		fmt.Println(i)
		for _, t := range triggers {
			if t.ID == i {
				id = i
				break
			}
		}
	}
	// it's not an int, so...
	for _, t := range triggers {
		if strings.ToLower(strings.Join(ctx.Args, " ")) == strings.ToLower(t.Trigger) {
			id = t.ID
			break
		}
	}

	err = ctx.Database.RemoveTrigger(ctx.Message.GuildID, id)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Sendf("Removed trigger %v.", id)
	return
}
