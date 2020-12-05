package bot

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/Starshine113/covebotnt/wlog"
	"github.com/bwmarrin/discordgo"
)

// Bot ...
type Bot struct {
	Session     *discordgo.Session
	Sugar       *wlog.Wlog
	Config      structs.BotConfig
	Pool        *cbdb.Db
	Bolt        *cbdb.BoltDb
	MemberCache *MemberCache
	UserCache   *UserCache
	RoleCache   *RoleCache
	Handlers    *ttlcache.Cache
	Version     string
	GitVer      string
	StartTime   time.Time
}

// NewBot returns a Bot struct
func NewBot(s *discordgo.Session, l *wlog.Wlog, p *cbdb.Db, b *cbdb.BoltDb, c structs.BotConfig, h *ttlcache.Cache, version, gitVer string, startTime time.Time) *Bot {
	bot := &Bot{Session: s, Sugar: l, Pool: p, Config: c, Bolt: b, Handlers: h, Version: version, GitVer: gitVer, StartTime: startTime}

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

	bot.Session.AddHandler(bot.GetGuildMembers)
	bot.Session.AddHandler(bot.GuildMemberChunk)
	bot.Session.AddHandler(bot.GuildMemberUpdate)
	bot.Session.AddHandler(bot.GuildRoleUpdate)

	return bot
}

// Prefix gets the prefix for the given guild
func (b *Bot) Prefix(guildID string) string {
	gs, err := b.Pool.GetGuildSettings(guildID)
	if err != nil {
		return b.Config.Bot.Prefixes[0]
	}

	if gs.Prefix != "" {
		return gs.Prefix
	}
	return b.Config.Bot.Prefixes[0]
}
