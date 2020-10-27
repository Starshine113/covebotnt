package main

import "github.com/bwmarrin/discordgo"

func checkOwner(userID string) error {
	for _, ownerID := range config.Bot.BotOwners {
		if userID == ownerID {
			return nil
		}
	}
	return &errorNoPermissions{"BotOwner"}
}

func checkAdmin(s *discordgo.Session, memberID, guildID string) error {
	// check if in DMs
	if guildID == "" {
		return &errorNoDMs{}
	}

	// get the guild
	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	// get the member
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// if not we check for admin perms
	// iterate through all guild roles
	for _, r := range guild.Roles {
		// iterate through member roles
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				// ...and the role has admin perms, return
				if r.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
					return nil
				}
			}
		}
	}
	return &errorNoPermissions{"Administrator"}
}

func checkModRole(s *discordgo.Session, memberID, guildID string, checkHelperRoles bool) error {
	// check if in DMs
	if guildID == "" {
		return &errorNoDMs{}
	}

	// get the guild
	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	// get the member
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// check if the user has a mod role
	for _, modRole := range globalSettings[guildID].Moderation.ModRoles {
		for _, role := range member.Roles {
			if role == modRole {
				return nil
			}
		}
	}

	// if this command only requires a helper role, check those too
	if checkHelperRoles {
		for _, helperRole := range globalSettings[guildID].Moderation.HelperRoles {
			for _, role := range member.Roles {
				if role == helperRole {
					return nil
				}
			}
		}
	}

	// if not we check for admin perms
	// iterate through all guild roles
	for _, r := range guild.Roles {
		// iterate through member roles
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				// ...and the role has admin perms, return
				if r.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
					return nil
				}
			}
		}
	}

	return &errorNoPermissions{"Administrator or ModRole"}
}
