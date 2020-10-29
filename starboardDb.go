package main

import (
	"context"
	"errors"
)

func getStarboardMessages() (map[string]string, map[string]string, error) {
	messageIDMap := make(map[string]string)
	starboardMsgIDMap := make(map[string]string)

	rows, err := db.Query(context.Background(), "select message_id, starboard_message_id from public.starboard_messages")
	if err != nil {
		return messageIDMap, starboardMsgIDMap, err
	}
	defer rows.Close()

	for rows.Next() {
		var messageID, starboardMsgID string
		rows.Scan(&messageID, &starboardMsgID)
		messageIDMap[messageID] = starboardMsgID
		starboardMsgIDMap[starboardMsgID] = messageID
	}

	return messageIDMap, starboardMsgIDMap, err
}

func insertStarboardEntry(messageID, channelID, guildID, starboardMessageID string) error {
	_, err := db.Exec(context.Background(), "insert into public.starboard_messages (message_id, channel_id, server_id, starboard_message_id) values ($1, $2, $3, $4)", messageID, channelID, guildID, starboardMessageID)
	if err != nil {
		return err
	}

	messageIDMap, starboardMsgIDMap, err = getStarboardMessages()
	return err
}

func deleteStarboardEntry(messageID string) error {
	sugar.Infof("Removing message entry for %v...", messageID)
	commandTag, err := db.Exec(context.Background(), "delete from public.starboard_messages where message_id = $1 or starboard_message_id = $2", messageID, messageID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}

	messageIDMap, starboardMsgIDMap, err = getStarboardMessages()
	return err
}

func setStarboardChannel(channelID, guildID string) error {
	sugar.Infof("Setting the starboard channel for %v to %v", guildID, channelID)
	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set starboard_channel = $1 where guild_id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	sugar.Infof("Set the starboard channel for %v to %v", guildID, channelID)
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}

func setStarboardLimit(limit int, guildID string) error {
	sugar.Infof("Setting the starboard limit for %v to %v", guildID, limit)
	commandTag, err := db.Exec(context.Background(), "update public.guild_settings set react_limit = $1 where guild_id = $2", limit, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	sugar.Infof("Set the starboard limit for %v to %v", guildID, limit)
	err = getSettingsForGuild(guildID)
	if err != nil {
		return err
	}
	sugar.Infof("Refreshed the settings for %v", guildID)
	return nil
}
