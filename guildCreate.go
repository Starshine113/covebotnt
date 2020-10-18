package main

import (
	"github.com/bwmarrin/discordgo"
)

func guildJoin(s *discordgo.Session, guild *discordgo.GuildCreate) {
	sugar.Debugf("Joined guild %v (%v)", guild.ID, guild.Name)
	// if this guild already has settings in the database, return
	if _, exists := globalSettings[guild.ID]; exists {
		return
	}

	err := initSettingsForGuild(guild.ID)
	if err != nil {
		sugar.Errorf("Error initialising the settings for guild %v: %v", guild.ID, err)
		return
	}
	sugar.Infof("Initialised settings for guild %v", guild.ID)
}
