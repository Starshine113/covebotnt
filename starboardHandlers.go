package main

import (
	"github.com/bwmarrin/discordgo"
)

func starboardReactionAdd(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	// if in DMs, do nothing
	if reaction.GuildID == "" {
		return
	}

	// if the channel is the server's starboard channel, do nothing
	if reaction.ChannelID == globalSettings[reaction.GuildID].Starboard.StarboardID {
		return
	}
	// if the emoji is not the server's starboard emoji, do nothing
	if reaction.Emoji.APIName() != globalSettings[reaction.GuildID].Starboard.Emoji {
		return
	}
	// if the channel is blacklisted, return
	for _, channel := range channelBlacklist[reaction.GuildID] {
		if channel == reaction.ChannelID {
			return
		}
	}

	// get the message
	message, err := s.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		sugar.Errorf("Error getting message %v: %v", reaction.MessageID, err)
	}

	// check the user, if it's the message author, remove the reaction
	if message.Author.ID == reaction.UserID {
		err = s.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.APIName(), reaction.UserID)
		if err != nil {
			sugar.Errorf("Error removing reaction on message %v: %v", reaction.MessageID, err)
		}
		return
	}

	// check reactions
	for _, messageReaction := range message.Reactions {
		if messageReaction.Emoji.APIName() == reaction.Emoji.APIName() {
			if messageReaction.Count >= globalSettings[reaction.GuildID].Starboard.ReactLimit {
				err = createOrEditMessage(s, message, reaction.GuildID, messageReaction.Count, messageReaction.Emoji)
				if err != nil {
					sugar.Errorf("Error creating/editing starboard message for %v: %v", message.ID, err)
				}
			} else if messageReaction.Count < globalSettings[reaction.GuildID].Starboard.ReactLimit {
				if _, ok := messageIDMap[message.ID]; ok {
					err := deleteStarboardMessage(s, message.ID, reaction.GuildID)
					if err != nil {
						sugar.Errorf("Error deleting message: %v", err)
					}
				}
			}
		}
	}
}

func starboardReactionRemove(s *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	// if the channel is the server's starboard channel, do nothing
	if reaction.ChannelID == globalSettings[reaction.GuildID].Starboard.StarboardID {
		return
	}
	// if the emoji is not the server's starboard emoji, do nothing
	if reaction.Emoji.APIName() != globalSettings[reaction.GuildID].Starboard.Emoji {
		return
	}
	// if the channel is blacklisted, return
	for _, channel := range channelBlacklist[reaction.GuildID] {
		if channel == reaction.ChannelID {
			return
		}
	}

	// get the message
	message, err := s.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		sugar.Errorf("Error getting message %v: %v", reaction.MessageID, err)
	}

	// check reactions
	for _, messageReaction := range message.Reactions {
		if messageReaction.Emoji.APIName() == reaction.Emoji.APIName() {
			if messageReaction.Count >= globalSettings[reaction.GuildID].Starboard.ReactLimit {
				err = createOrEditMessage(s, message, reaction.GuildID, messageReaction.Count, messageReaction.Emoji)
				if err != nil {
					sugar.Errorf("Error creating/editing starboard message for %v: %v", message.ID, err)
				}
			} else if messageReaction.Count < globalSettings[reaction.GuildID].Starboard.ReactLimit {
				if _, ok := messageIDMap[message.ID]; ok {
					err := deleteStarboardMessage(s, messageIDMap[message.ID], reaction.GuildID)
					if err != nil {
						sugar.Errorf("Error deleting message: %v", err)
					}
				}
			}
		}
	}
}

func starboardMessageDelete(s *discordgo.Session, message *discordgo.MessageDelete) {
	if _, ok := messageIDMap[message.ID]; ok {
		err := deleteStarboardEntry(messageIDMap[message.ID])
		if err != nil {
			sugar.Errorf("Error deleting message entry for %v: %v", message.ID, err)
			return
		}
		sugar.Infof("Deleted starboard database entry for %v", message.ID)
	}
	if _, ok := starboardMsgIDMap[message.ID]; ok {
		err := deleteStarboardEntry(starboardMsgIDMap[message.ID])
		if err != nil {
			sugar.Errorf("Error deleting message entry for %v: %v", message.ID, err)
			return
		}
		sugar.Infof("Deleted starboard database entry for %v", message.ID)
	}
}
