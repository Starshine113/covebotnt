package main

import (
	"os"

	"github.com/Starshine113/covebotnt/crouter"
)

func commandKill(ctx *crouter.Ctx) (err error) {
	_, err = ctx.Send("Restarting the bot, please wait...")
	if err != nil {
		return err
	}
	sugar.Infof("Kill command received, shutting down...")

	dg.Close()
	sugar.Infof("Disconnected from Discord.")
	db.Close()
	sugar.Infof("Closed database connection.")

	logger.Sync()
	boltDb.Bolt.Close()
	levelCache.Close()

	os.Exit(0)
	return nil
}
