package commands

import (
	"fmt"
	"time"

	"github.com/Starshine113/covebotnt/cbctx"
)

// Ping command: replies with latency and message edit time
func Ping(ctx *cbctx.Ctx) (err error) {
	heartbeat := ctx.Session.HeartbeatLatency().String()

	// get current time
	cmdStart := time.Now()

	// send initial message
	message, err := ctx.Send("Pong!\nHeartbeat: " + heartbeat)
	if err != nil {
		return fmt.Errorf("Ping: %w", err)
	}

	// get time difference, edit message
	diff := time.Now().Sub(cmdStart).String()
	_, err = ctx.Edit(message, message.Content+"\nMessage latency: "+diff)
	return err
}
