package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Get ...
func (c *MemberCache) Get(gID, mID string) (m *discordgo.Member, err error) {
	key := fmt.Sprintf("%v-%v", gID, mID)
	c.Cache.mu.Lock()
	if v, e := c.Cache.Cache.Get(key); e == nil {
		c.Cache.mu.Unlock()
		return v.(*discordgo.Member), nil
	}

	m, err = c.Cache.s.GuildMember(gID, mID)
	if err != nil {
		c.Cache.mu.Unlock()
		return nil, err
	}
	c.Cache.Cache.Set(key, m)
	c.Cache.mu.Unlock()
	return m, nil
}

// Remove ...
func (c *MemberCache) Remove(gID, mID string) (err error) {
	key := fmt.Sprintf("%v-%v", gID, mID)
	c.Cache.mu.Lock()
	err = c.Cache.Cache.Remove(key)
	c.Cache.mu.Unlock()
	return err
}

// ForceGet forces a cache refresh for a member
func (c *MemberCache) ForceGet(gID, mID string) (m *discordgo.Member, err error) {
	c.Remove(gID, mID)
	return c.Get(gID, mID)
}

// AddNoExpire adds a member to the cache, with no expire time
func (c *MemberCache) AddNoExpire(gID, mID string, m *discordgo.Member) {
	key := fmt.Sprintf("%v-%v", gID, mID)
	c.Cache.mu.Lock()
	c.Cache.Cache.SetWithTTL(key, m, time.Duration(0))
	c.Cache.mu.Unlock()
}

// Add manually adds a member to the cache
func (c *MemberCache) Add(gID, mID string, m *discordgo.Member) {
	key := fmt.Sprintf("%v-%v", gID, mID)
	c.Cache.mu.Lock()
	c.Cache.Cache.Set(key, m)
	c.Cache.mu.Unlock()
}
