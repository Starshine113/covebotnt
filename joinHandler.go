package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func onJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	fmt.Printf("New member joined %v: %v#%v (%v)\n", member.GuildID, member.User.Username, member.User.Discriminator, member.User.ID)
}
