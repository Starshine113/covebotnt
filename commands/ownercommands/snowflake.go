package ownercommands

import (
	"github.com/starshine-sys/covebotnt/crouter"
)

func snowflake(ctx *crouter.Ctx) (err error) {
	s := ctx.Bot.SnowflakeGen.Get()
	t, err := ctx.Bot.SnowflakeGen.Parse(s)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	_, err = ctx.EmbedfNoXHandler("Generated ID", "```%v```\n%v", s, t.Format("Jan _2 2006, 15:04:05 MST"))
	return
}
