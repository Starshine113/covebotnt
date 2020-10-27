package main

import (
	"github.com/bwmarrin/discordgo"
)

func guildJoin(s *discordgo.Session, guild *discordgo.GuildCreate) {
	var err error

	sugar.Debugf("Joined guild %v (%v)", guild.ID, guild.Name)
	// init database for this guild
	err = boltDb.InitForGuild(guild.ID)
	if err != nil {
		sugar.Errorf("Error initialising cache db for guild %v: %v", guild.ID, err)
	}

	// if this guild already has settings in the database, return
	if _, exists := globalSettings[guild.ID]; exists {
		return
	}

	err = initSettingsForGuild(guild.ID)
	if err != nil {
		sugar.Errorf("Error initialising the settings for guild %v: %v", guild.ID, err)
		return
	}
	sugar.Infof("Initialised settings for guild %v", guild.ID)
}
