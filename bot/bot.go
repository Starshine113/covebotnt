package bot

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/bwmarrin/discordgo"
)

// Bot ...
type Bot struct {
	Session     *discordgo.Session
	MemberCache *MemberCache
	UserCache   *UserCache
	RoleCache   *RoleCache
}

// NewBot returns a Bot struct
func NewBot(s *discordgo.Session) *Bot {
	bot := &Bot{Session: s}

	m := ttlcache.NewCache()
	m.SetTTL(time.Hour)
	m.SkipTTLExtensionOnHit(true)
	bot.MemberCache = &MemberCache{&Cache{Cache: m, s: s}}

	u := ttlcache.NewCache()
	u.SetTTL(time.Hour)
	m.SkipTTLExtensionOnHit(true)
	bot.UserCache = &UserCache{&Cache{Cache: u, s: s}}

	r := ttlcache.NewCache()
	u.SetTTL(time.Hour)
	m.SkipTTLExtensionOnHit(true)
	bot.RoleCache = &RoleCache{&Cache{Cache: r, s: s}}

	return bot
}
