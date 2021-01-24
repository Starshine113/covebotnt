package admincommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// HelperRoles ...
func HelperRoles(ctx *crouter.Ctx) (err error) {
	if err = ctx.TriggerTyping(); err != nil {
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
			err = addHelperRole(ctx, ctx.Message.GuildID, role.ID)
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
			err = delHelperRole(ctx, ctx.Message.GuildID, role.ID)
			if err != nil {
				ctx.CommandError(err)
				return nil
			}
			_, err = ctx.Send(fmt.Sprintf("Removed role **%v** from the list of helper roles.", role.Name))
			return err
		}
	}
	roleEmbed, err := getHelperRoles(ctx, ctx.Message.GuildID)
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

func addHelperRole(ctx *crouter.Ctx, guildID, roleID string) (err error) {
	gs, err := ctx.Database.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	roles := append(gs.Moderation.HelperRoles, roleID)

	commandTag, err := ctx.Database.Pool.Exec(context.Background(), "update public.guild_settings set helper_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = ctx.Database.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	ctx.Bot.Sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func delHelperRole(ctx *crouter.Ctx, guildID, roleID string) (err error) {
	gs, err := ctx.Database.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	roles := gs.Moderation.HelperRoles
	for i, v := range roles {
		if v == roleID {
			roles, err = removeElement(roles, i)
			if err != nil {
				return
			}
		}
	}

	commandTag, err := ctx.Database.Pool.Exec(context.Background(), "update public.guild_settings set helper_roles = $1 where guild_id = $2", roles, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	err = ctx.Database.RemoveFromGuildCache(guildID)
	if err != nil {
		return err
	}
	ctx.Bot.Sugar.Infof("Refreshed the settings for %v", guildID)
	return
}

func getHelperRoles(ctx *crouter.Ctx, guildID string) (*discordgo.MessageEmbed, error) {
	var roles []string
	guild, err := ctx.Session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	gs, err := ctx.Database.GetGuildSettings(guildID)
	if err != nil {
		return nil, err
	}
	for _, modRole := range gs.Moderation.HelperRoles {
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
