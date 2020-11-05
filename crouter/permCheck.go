package crouter

import (
	"github.com/Starshine113/covebotnt/cbctx"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

func checkHelperPerm(ctx *cbctx.Ctx, guildSettings *structs.GuildSettings) (err error) {
	// check if in DMs
	if ctx.Message.GuildID == "" {
		return &cbctx.ErrorNoDMs{}
	}

	// get the guild
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == discordgo.ErrStateNotFound {
		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// get the member
	member, err := ctx.Session.State.Member(ctx.Message.GuildID, ctx.Author.ID)
	if err == discordgo.ErrStateNotFound {
		member, err = ctx.Session.GuildMember(ctx.Message.GuildID, ctx.Author.ID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// check if the user has a mod role
	for _, modRole := range guildSettings.Moderation.ModRoles {
		for _, role := range member.Roles {
			if role == modRole {
				return nil
			}
		}
	}

	// if this command only requires a helper role, check those too
	for _, helperRole := range guildSettings.Moderation.HelperRoles {
		for _, role := range member.Roles {
			if role == helperRole {
				return nil
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

	return &cbctx.ErrorNoPermissions{MissingPerms: "Administrator, HelperRole, or ModRole"}
}

func checkModPerm(ctx *cbctx.Ctx, guildSettings *structs.GuildSettings) (err error) {
	// check if in DMs
	if ctx.Message.GuildID == "" {
		return &cbctx.ErrorNoDMs{}
	}

	// get the guild
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == discordgo.ErrStateNotFound {
		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// get the member
	member, err := ctx.Session.State.Member(ctx.Message.GuildID, ctx.Author.ID)
	if err == discordgo.ErrStateNotFound {
		member, err = ctx.Session.GuildMember(ctx.Message.GuildID, ctx.Author.ID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// check if the user has a mod role
	for _, modRole := range guildSettings.Moderation.ModRoles {
		for _, role := range member.Roles {
			if role == modRole {
				return nil
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

	return &cbctx.ErrorNoPermissions{MissingPerms: "Administrator or ModRole"}
}

func checkAdmin(ctx *cbctx.Ctx) (err error) {
	// check if in DMs
	if ctx.Message.GuildID == "" {
		return &cbctx.ErrorNoDMs{}
	}

	// get the guild
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == discordgo.ErrStateNotFound {
		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// get the member
	member, err := ctx.Session.State.Member(ctx.Message.GuildID, ctx.Author.ID)
	if err == discordgo.ErrStateNotFound {
		member, err = ctx.Session.GuildMember(ctx.Message.GuildID, ctx.Author.ID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
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

	return &cbctx.ErrorNoPermissions{MissingPerms: "Administrator"}
}

func checkOwner(userID string, owners []string) error {
	for _, ownerID := range owners {
		if userID == ownerID {
			return nil
		}
	}
	return &cbctx.ErrorNoPermissions{MissingPerms: "BotOwner"}
}
