package etc

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// PermError is a permission error
type PermError struct {
	perms int64
	s     []string
}

func (p *PermError) Error() string {
	if len(p.s) > 0 {
		return strings.Join(p.s, ", ")
	}
	return strconv.FormatInt(p.perms, 10)
}

var permNames = map[int64]string{
	discordgo.PermissionCreateInstantInvite:  "Create Instant Invite",
	discordgo.PermissionKickMembers:          "Kick Members",
	discordgo.PermissionBanMembers:           "Ban Members",
	discordgo.PermissionAdministrator:        "Administrator",
	discordgo.PermissionManageChannels:       "Manage Channels",
	discordgo.PermissionManageServer:         "Manage Server",
	discordgo.PermissionAddReactions:         "Add Reactions",
	discordgo.PermissionViewAuditLogs:        "View Audit Log",
	discordgo.PermissionVoicePrioritySpeaker: "Priority Speaker",
	discordgo.PermissionViewChannel:          "View Channel",
	discordgo.PermissionSendMessages:         "Send Messages",
	discordgo.PermissionSendTTSMessages:      "Send TTS Messages",
	discordgo.PermissionManageMessages:       "Manage Messages",
	discordgo.PermissionEmbedLinks:           "Embed Links",
	discordgo.PermissionAttachFiles:          "Attach Files",
	discordgo.PermissionReadMessageHistory:   "Read Message History",
	discordgo.PermissionMentionEveryone:      "Mention Everyone",
	discordgo.PermissionUseExternalEmojis:    "Use External Emojis",
	discordgo.PermissionVoiceConnect:         "Connect",
	discordgo.PermissionVoiceSpeak:           "Speak",
	discordgo.PermissionVoiceMuteMembers:     "Mute Members",
	discordgo.PermissionVoiceDeafenMembers:   "Deafen Members",
	discordgo.PermissionVoiceMoveMembers:     "Move Members",
	discordgo.PermissionVoiceUseVAD:          "Use VAD",
	discordgo.PermissionChangeNickname:       "Change Nickname",
	discordgo.PermissionManageNicknames:      "Manage Nicknames",
	discordgo.PermissionManageRoles:          "Manage Roles",
	discordgo.PermissionManageWebhooks:       "Manage Webhooks",
	discordgo.PermissionManageEmojis:         "Manage Emojis",
}

// PermStrings gets the permission strings for all required permissions
func PermStrings(p int64) []string {
	var out = make([]string, 0, 32)
	for perm, name := range permNames {
		if p&perm == perm {
			out = append(out, name)
		}
	}

	return out
}
