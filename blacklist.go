package main

import (
	"context"
	"fmt"
)

func getBlacklistAll() map[uint64][]uint64 {
	channelBlacklist := make(map[uint64][]uint64)
	var guilds []uint64

	guildRows, err := db.Query(context.Background(), "select distinct server_id from starboard_blacklisted_channels")
	if err != nil {
		panic(err)
	}
	defer guildRows.Close()

	for guildRows.Next() {
		var guild uint64
		err = guildRows.Scan(&guild)
		guilds = append(guilds, guild)
	}

	for _, guild := range guilds {
		var channels []uint64

		channelRows, err := db.Query(context.Background(), "select channel_id from starboard_blacklisted_channels where server_id=$1", guild)
		if err != nil {
			panic(err)
		}
		defer channelRows.Close()

		for channelRows.Next() {
			var channel uint64
			err = channelRows.Scan(&channel)
			if err != nil {
				panic(err)
			}
			channels = append(channels, channel)
		}

		channelBlacklist[guild] = channels
	}
	return channelBlacklist
}

func getBlacklistForGuild(guildID uint64) error {
	var channels []uint64

	channelRows, err := db.Query(context.Background(), "select channel_id from starboard_blacklisted_channels where server_id=$1", guildID)
	if err != nil {
		return err
	}
	defer channelRows.Close()

	for channelRows.Next() {
		var channel uint64
		err = channelRows.Scan(&channel)
		if err != nil {
			return err
		}
		channels = append(channels, channel)
	}

	channelBlacklist[guildID] = channels
	return nil
}

func addChannelsToBlacklist(guildID uint64, channels []uint64) error {
	for _, channel := range channelBlacklist[guildID] {
		for i, newChannel := range channels {
			if channel == newChannel {
				channels = removeFromSlice(channels, i)
			}
		}
	}

	for _, channel := range channels {
		commandTag, err := db.Exec(context.Background(), "insert into starboard_blacklisted_channels (channel_id, server_id) values ($1, $2)", channel, guildID)
		if err != nil {
			return err
		}
		if commandTag.RowsAffected() == 0 {
			return fmt.Errorf("Channel %v (Guilld %v) was not added to the blacklist", channel, guildID)
		}
	}

	err := getBlacklistForGuild(guildID)
	if err != nil {
		return err
	}
	return nil
}

func removeFromSlice(s []uint64, i int) []uint64 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
