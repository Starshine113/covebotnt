package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Snowflake shows the timestamp of all discord IDs given
func Snowflake(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		ctx.CommandError(err)
		return nil
	}

	var msgs []string
	for _, arg := range ctx.Args {
		arg = strings.TrimSpace(arg)
		args := strings.Split(arg, "\n")
		for _, a := range args {
			t, err := discordgo.SnowflakeTimestamp(a)
			if err != nil {
				ctx.CommandError(err)
				continue
			}
			msgs = append(msgs, fmt.Sprintf("`%v`: %v", a, t.UTC().Format(time.RFC3339)))
		}
	}

	desc := strings.Join(msgs, "\n")
	if len(desc) > 2000 {
		desc = desc[:2000]
	}
	_, err = ctx.Send(&discordgo.MessageEmbed{
		Title:       "Timestamps",
		Description: desc,
	})
	return
}
