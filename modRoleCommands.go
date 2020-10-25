package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/bwmarrin/discordgo"
)

func commandModRoles(ctx *cbctx.Ctx) (err error) {
	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	perms := checkAdmin(ctx.Session, ctx.Message.Author.ID, ctx.Message.GuildID)
	if perms != nil {
		commandError(perms, ctx.Session, ctx.Message)
		return nil
	}

	if len(ctx.Args) >= 1 {
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) != 2 {
				_, err = ctx.Send(cbctx.ErrorEmoji + "No/too many roles supplied.")
				return err
			}
			err = addModRole(ctx.Message.GuildID, ctx.Args[1])
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Added role **%v** to the list of moderator roles.", ctx.Args[0]))
			return err
		case "remove", "delete":
			_, err = ctx.Send("remove a role")
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

func addModRole(guildID, roleID string) (err error) {
	roles := append(globalSettings[guildID].Moderation.ModRoles, roleID)

	commandTag, err := db.Exec(context.Background(), "update guild_settings set mod_roles = $1 where guild_id = $2", roles, guildID)
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

func addHelperRole(guildID, roleID string) (err error) {
	roles := append(globalSettings[guildID].Moderation.HelperRoles, roleID)

	commandTag, err := db.Exec(context.Background(), "update guild_settings helper_roles = $1 where guild_id = $2", roles, guildID)
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
