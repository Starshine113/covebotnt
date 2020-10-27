package levels

import (
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/bwmarrin/discordgo"
)

// MessageCreate is the message create handler for XP
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, cache *ttlcache.Cache, db *cbdb.BoltDb) {
	if _, exists := cache.Get(m.Author.ID); exists == nil {
		db.AddXPForUser(m.Author.ID, m.GuildID)
		return
	}
	cache.Set(m.Author.ID, true)
	return
}
