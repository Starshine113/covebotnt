package main

import (
	"os"

	"github.com/Starshine113/covebotnt/cbctx"
)

func commandKill(ctx *cbctx.Ctx) (err error) {
	_, err = ctx.Send("<:blobsob:766276787814531093> Shutting down...")
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
