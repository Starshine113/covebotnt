package ownercommands

import (
	"context"

	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/jackc/pgx/v4"
)

func fixSnowflakes(ctx *crouter.Ctx) (err error) {
	ids := make([]int, 0)

	err = ctx.Database.Pool.QueryRow(context.Background(), "select array(select id from public.mod_log where snowflake is null)").Scan(&ids)
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	b := &pgx.Batch{}

	for _, id := range ids {
		s := ctx.Bot.SnowflakeGen.Get()
		b.Queue("update public.mod_log set snowflake = $1 where id = $2", s, id)
	}

	results := ctx.Database.Pool.SendBatch(context.Background(), b)
	defer results.Close()

	_, err = ctx.Sendf("Added snowflakes to %v entries.", len(ids))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	ids = make([]int, 0)

	err = ctx.Database.Pool.QueryRow(context.Background(), "select array(select id from public.triggers where snowflake is null)").Scan(&ids)
	if err != nil {
		_, err = ctx.CommandError(err)
		return
	}

	b = &pgx.Batch{}

	for _, id := range ids {
		s := ctx.Bot.SnowflakeGen.Get()
		b.Queue("update public.triggers set snowflake = $1 where id = $2", s, id)
	}

	results = ctx.Database.Pool.SendBatch(context.Background(), b)
	defer results.Close()

	_, err = ctx.Sendf("Added snowflakes to %v triggers.", len(ids))
	return err
}
