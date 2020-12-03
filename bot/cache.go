package bot

import (
	"sync"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/bwmarrin/discordgo"
)

// Cache ...
type Cache struct {
	s     *discordgo.Session
	Cache *ttlcache.Cache
	mu    sync.Mutex
	t     cacheType
}

// MemberCache ...
type MemberCache struct {
	*Cache
}

// UserCache ...
type UserCache struct {
	*Cache
}

// RoleCache ...
type RoleCache struct {
	*Cache
}

type cacheType int

const (
	memberCache cacheType = iota
	userCache
	roleCache
)
