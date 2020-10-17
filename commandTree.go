package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandTree(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch command {
	case "ping":
		commandPing(s, m)
	case "help":
		commandHelp(args, s, m)
	case "setstatus":
		commandSetStatus(args, s, m)
	case "starboard":
		commandStarboard(args, s, m)
	}
}
