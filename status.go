package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	go updateStatusLoop(s)
}

func updateStatusLoop(s *discordgo.Session) {
	for {
		newStatus := config.Bot.Prefixes[0] + "help | in " + fmt.Sprint(len(s.State.Guilds)) + " guilds"
		if config.Bot.CustomStatus.Status != "" {
			newStatus += " | " + config.Bot.CustomStatus.Status
		}
		if config.Bot.CustomStatus.Override {
			newStatus = config.Bot.CustomStatus.Status
		}
		err := dg.UpdateStatus(0, newStatus)
		if err != nil {
			sugar.Errorf("Update status error: ", err)
		}
		time.Sleep(time.Minute)
	}
}
