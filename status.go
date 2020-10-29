package main

import (
	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	err := updateStatus(s)
	if err != nil {
		sugar.Errorf("Error updating status: %v", err)
	}
}

func updateStatus(s *discordgo.Session) (err error) {
	switch config.Bot.CustomStatus.Type {
	case "listening":
		err = s.UpdateListeningStatus(config.Bot.CustomStatus.Status)
	case "playing":
		err = s.UpdateStatus(0, config.Bot.CustomStatus.Status)
	default:
		return nil
	}
	return err
}
