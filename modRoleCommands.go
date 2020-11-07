package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

func commandModRoles(ctx *crouter.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	if len(ctx.Args) >= 1 {
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) != 2 {
				_, err = ctx.Send(crouter.ErrorEmoji + "No/too many roles supplied.")
				return err
			}
			role, err := ctx.ParseRole(ctx.Args[1])
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			err = addModRole(ctx.Message.GuildID, role.ID)
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Added role **%v** to the list of moderator roles.", role.Name))
			return err
		case "remove", "delete":
			if len(ctx.Args) != 2 {
				_, err = ctx.Send(crouter.ErrorEmoji + "No/too many roles supplied.")
				return err
			}
			role, err := ctx.ParseRole(ctx.Args[1])
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			err = delModRole(ctx.Message.GuildID, role.ID)
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Removed role **%v** from the list of moderator roles.", role.Name))
			return err
		}
	}
	roleEmbed, err := getModRoles(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send(roleEmbed)
	if err != nil {
		return fmt.Errorf("ModRoles: %v", err)
	}
	return nil
}

func commandHelperRoles(ctx *crouter.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	if len(ctx.Args) >= 1 {
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) != 2 {
				_, err = ctx.Send(crouter.ErrorEmoji + "No/too many roles supplied.")
				return err
			}
			role, err := ctx.ParseRole(ctx.Args[1])
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			err = addHelperRole(ctx.Message.GuildID, role.ID)
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Added role **%v** to the list of helper roles.", role.Name))
			return err
		case "remove", "delete":
			if len(ctx.Args) != 2 {
				_, err = ctx.Send(crouter.ErrorEmoji + "No/too many roles supplied.")
				return err
			}
			role, err := ctx.ParseRole(ctx.Args[1])
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			err = delHelperRole(ctx.Message.GuildID, role.ID)
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Removed role **%v** from the list of helper roles.", role.Name))
			return err
		}
	}
	roleEmbed, err := getHelperRoles(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send(roleEmbed)
	if err != nil {
		return fmt.Errorf("ModRoles: %v", err)
	}
	return nil
}

func addModRole(guildID, roleID string) (err error) {
	roles := append(globalSettings[guildID].Moderation.ModRoles, roleID)

	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set mod_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func delModRole(guildID, roleID string) (err error) {
	roles := globalSettings[guildID].Moderation.ModRoles
	for i, v := range roles {
		if v == roleID {
			roles, err = removeElement(roles, i)
			if err != nil {
				return
			}
		}
	}

	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set mod_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return
}

func addHelperRole(guildID, roleID string) (err error) {
	roles := append(globalSettings[guildID].Moderation.HelperRoles, roleID)

	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set helper_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func delHelperRole(guildID, roleID string) (err error) {
	roles := globalSettings[guildID].Moderation.HelperRoles
	for i, v := range roles {
		if v == roleID {
			roles, err = removeElement(roles, i)
			if err != nil {
				return
			}
		}
	}

	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set helper_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return
}

func getModRoles(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	var roles []string
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, modRole := range globalSettings[guildID].Moderation.ModRoles {
		roles = append(roles, "<@&"+modRole+">")
	}

	var roleString string
	if len(roles) == 0 {
		roleString = "No mod roles."
	} else {
		for i, role := range roles {
			if i == len(roles)-1 {
				roleString += role
			} else {
				roleString += role + ", "
			}
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Mod roles for " + guild.Name,
		Description: roleString,
		Color:       0x21a1a8,
	}
	return embed, nil
}

func getHelperRoles(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	var roles []string
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, modRole := range globalSettings[guildID].Moderation.HelperRoles {
		roles = append(roles, "<@&"+modRole+">")
	}

	var roleString string
	if len(roles) == 0 {
		roleString = "No helper roles."
	} else {
		for i, role := range roles {
			if i == len(roles)-1 {
				roleString += role
			} else {
				roleString += role + ", "
			}
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Helper roles for " + guild.Name,
		Description: roleString,
		Color:       0x21a1a8,
	}
	return embed, nil
}

func removeElement(s []string, i int) ([]string, error) {
	// s is [1,2,3,4,5,6], i is 2

	// perform bounds checking first to prevent a panic!
	if i >= len(s) || i < 0 {
		return nil, fmt.Errorf("Index is out of range. Index is %d with slice length %d", i, len(s))
	}

	// copy the last element (6) to index `i`. At this point,
	// `s` will be [1,2,6,4,5,6]
	s[i] = s[len(s)-1]
	// Remove the last element from the slice by truncating it
	// This way, `s` will now include all the elements from index 0
	// up to (but not including) the last element
	return s[:len(s)-1], nil
}
