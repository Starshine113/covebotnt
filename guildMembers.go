package main

import "github.com/bwmarrin/discordgo"

func getGuildMembers(s *discordgo.Session, _ *discordgo.Ready) {
	for _, g := range s.State.Guilds {
		s.RequestGuildMembers(g.ID, "", 0, false)
	}
}

func guildMemberChunk(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
	sugar.Debugf("Received %v members for guild %v", len(chunk.Members), chunk.GuildID)
	for _, m := range chunk.Members {
		b.MemberCache.AddNoExpire(m.GuildID, m.User.ID, m)
	}
}

func guildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	sugar.Debugf("Received updated data for %v in %v", m.User.ID, m.GuildID)
	b.MemberCache.Remove(m.GuildID, m.User.ID)
	b.MemberCache.AddNoExpire(m.GuildID, m.User.ID, m.Member)
}
