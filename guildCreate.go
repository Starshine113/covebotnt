package main

import (
	"github.com/bwmarrin/discordgo"
)

func guildJoin(s *discordgo.Session, guild *discordgo.GuildCreate) {
	var err error

	sugar.Debugf("Joined guild %v (%v)", guild.ID, guild.Name)

	err = pool.InitSettingsForGuild(guild.ID)
	if err != nil {
		sugar.Errorf("Error initialising the settings for guild %v: %v", guild.ID, err)
		return
	}
	sugar.Infof("Initialised settings for guild %v", guild.ID)
}
