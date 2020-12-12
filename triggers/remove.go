package triggers

import (
	"strings"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/snowflake"
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
	var s snowflake.Snowflake

	// check if it's a snowflake
	if cbdb.Snowflake.MatchString(ctx.Args[0]) {
		for _, t := range triggers {
			if string(t.Snowflake) == ctx.Args[0] {
				id = t.ID
				s = t.Snowflake
				break
			}
		}
	} else {
		// it's not a snowflake, so...
		for _, t := range triggers {
			if strings.ToLower(strings.Join(ctx.Args, " ")) == strings.ToLower(t.Trigger) {
				id = t.ID
				s = t.Snowflake
				break
			}
		}
	}

	err = ctx.Database.RemoveTrigger(ctx.Message.GuildID, id)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	_, err = ctx.Sendf("Removed trigger `%v`.", s)
	return
}
