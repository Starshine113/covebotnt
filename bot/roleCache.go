package bot

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Get ...
func (c *RoleCache) Get(gID, rID string) (r *discordgo.Role, err error) {
	key := fmt.Sprintf("%v-%v", gID, rID)
	c.Cache.mu.Lock()
	if v, e := c.Cache.Cache.Get(key); e == nil {
		c.Cache.mu.Unlock()
		return v.(*discordgo.Role), nil
	}

	roles, err := c.Cache.s.GuildRoles(gID)
	if err != nil {
		c.Cache.mu.Unlock()
		return nil, err
	}
	for _, role := range roles {
		c.Cache.Cache.Set(fmt.Sprintf("%v-%v", gID, role.ID), role)
		if role.ID == rID {
			r = role
		}
	}
	c.Cache.mu.Unlock()
	if r == nil {
		return nil, errors.New("role not found")
	}
	return r, nil
}

// Remove ...
func (c *RoleCache) Remove(gID, rID string) (err error) {
	key := fmt.Sprintf("%v-%v", gID, rID)
	c.Cache.mu.Lock()
	err = c.Cache.Cache.Remove(key)
	c.Cache.mu.Unlock()
	return err
}

// ForceGet forces a cache refresh for a member
func (c *RoleCache) ForceGet(gID, mID string) (r *discordgo.Role, err error) {
	c.Remove(gID, mID)
	return c.Get(gID, mID)
}

// AddNoExpire adds a member to the cache, with no expire time
func (c *RoleCache) AddNoExpire(gID, rID string, r *discordgo.Role) {
	key := fmt.Sprintf("%v-%v", gID, rID)
	c.Cache.mu.Lock()
	c.Cache.Cache.SetWithTTL(key, rID, time.Duration(0))
	c.Cache.mu.Unlock()
}

// Add manually adds a member to the cache
func (c *RoleCache) Add(gID, mID string, m *discordgo.Role) {
	key := fmt.Sprintf("%v-%v", gID, mID)
	c.Cache.mu.Lock()
	c.Cache.Cache.Set(key, m)
	c.Cache.mu.Unlock()
}
