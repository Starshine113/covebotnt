package usercommands

import (
	"github.com/Starshine113/covebotnt/crouter"
	d "github.com/bwmarrin/discordgo"
)

// Invite sends an invite link for the bot
func Invite(ctx *crouter.Ctx) (err error) {
	// perms is the list of permissions the bot will be granted by default
	var perms = d.PermissionReadMessages +
		d.PermissionReadMessageHistory +
		d.PermissionSendMessages +
		d.PermissionManageMessages +
		d.PermissionManageEmojis +
		d.PermissionChangeNickname +
		d.PermissionEmbedLinks +
		d.PermissionAttachFiles +
		d.PermissionUseExternalEmojis +
		d.PermissionAddReactions

	_, err = ctx.Sendf("Use this link to invite me to your server: <https://discord.com/oauth2/authorize?client_id=%v&scope=bot&permissions=%v>", ctx.BotUser.ID, perms)
	return
}
