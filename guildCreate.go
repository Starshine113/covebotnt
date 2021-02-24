package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func guildJoin(s *discordgo.Session, guild *discordgo.GuildCreate) {
	var err error

	sugar.Debugf("Joined guild %v (%v)", guild.ID, guild.Name)

	// leave banned guilds
	for _, id := range config.Bot.BannedServers {
		if id == guild.ID {
			err = s.GuildLeave(guild.ID)
			if err != nil {
				sugar.Errorf("Error leaving guild %v (%v): %v", guild.Name, guild.ID, err)
				return
			}
			sugar.Infof("Automatically left banned guild %v (%v).", guild.Name, guild.ID)
			return
		}
	}
	// if the owner is a blocked user, leave the guild too
	for _, id := range config.Bot.BlockedUsers {
		if id == guild.OwnerID {
			err = s.GuildLeave(guild.ID)
			if err != nil {
				sugar.Errorf("Error leaving guild %v (%v): %v", guild.Name, guild.ID, err)
				return
			}
			sugar.Infof("Automatically left guild %v (%v), owned by blocked user %v.", guild.Name, guild.ID, guild.OwnerID)
			return
		}
	}

	err = pool.InitSettingsForGuild(guild.ID)
	if err != nil {
		sugar.Errorf("Error initialising the settings for guild %v: %v", guild.ID, err)
		return
	}
	sugar.Infof("Initialised settings for guild %v", guild.ID)

	for _, r := range guild.Roles {
		b.RoleCache.Cache.Cache.SetWithTTL(fmt.Sprintf("%v-%v", guild.ID, r.ID), r, 0)
	}
}
