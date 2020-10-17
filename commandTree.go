package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandTree(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if command == "ping" {
		commandPing(s, m)
	} else if command == "help" {
		commandHelp(args, s, m)
	} else if command == "setstatus" {
		commandSetStatus(args, s, m)
	}
}
