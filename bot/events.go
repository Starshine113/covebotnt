package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// GetGuildMembers ...
func (b *Bot) GetGuildMembers(s *discordgo.Session, _ *discordgo.Ready) {
	for _, g := range s.State.Guilds {
		s.RequestGuildMembers(g.ID, "", 0, false)
	}
}

// GuildMemberChunk ...
func (b *Bot) GuildMemberChunk(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
	b.Sugar.Debugf("Received %v members for guild %v", len(chunk.Members), chunk.GuildID)
	for _, m := range chunk.Members {
		b.MemberCache.AddNoExpire(m.GuildID, m.User.ID, m)
	}
}

// GuildMemberUpdate ...
func (b *Bot) GuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	b.Sugar.Debugf("Received updated data for %v in %v", m.User.ID, m.GuildID)
	b.MemberCache.Remove(m.GuildID, m.User.ID)
	b.MemberCache.AddNoExpire(m.GuildID, m.User.ID, m.Member)
}

// GuildRoleUpdate ...
func (b *Bot) GuildRoleUpdate(s *discordgo.Session, r *discordgo.GuildRoleUpdate) {
	b.Sugar.Debugf("Received updated data for %v in %v", r.Role.ID, r.GuildID)
	b.RoleCache.Remove(r.GuildID, r.Role.ID)
	b.RoleCache.Cache.Cache.SetWithTTL(fmt.Sprintf("%v-%v", r.GuildID, r.Role.ID), r.Role, 0)
}
