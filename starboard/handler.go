package starboard

import (
	"log"
	"sync"

	"github.com/starshine-sys/covebotnt/cbdb"
	"github.com/starshine-sys/covebotnt/wlog"
	"github.com/bwmarrin/discordgo"
)

// Sb ...
type Sb struct {
	Sugar *wlog.Wlog
	Pool  *cbdb.Db
	mu    sync.Mutex
}

// ReactionAdd ...
func (sb *Sb) ReactionAdd(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	// if in DMs, do nothing
	if reaction.GuildID == "" {
		return
	}

	gs, err := sb.Pool.GetGuildSettings(reaction.GuildID)
	if err != nil {
		return
	}

	// if the channel is the server's starboard channel, do nothing
	if reaction.ChannelID == gs.Starboard.StarboardID {
		return
	}
	// if the emoji is not the server's starboard emoji, do nothing
	if reaction.Emoji.MessageFormat() != gs.Starboard.Emoji {
		return
	}

	if sb.Pool.InStarboardBlacklist(reaction.GuildID, reaction.ChannelID) {
		return
	}

	// get the message
	message, err := s.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		sb.Sugar.Errorf("Error getting message %v: %v", reaction.MessageID, err)
	}

	// check the user, if it's the message author, remove the reaction
	if message.Author.ID == reaction.UserID && !gs.Starboard.SenderCanReact {
		err = s.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.APIName(), reaction.UserID)
		if err != nil {
			sb.Sugar.Errorf("Error removing reaction on message %v: %v", reaction.MessageID, err)
		}
		return
	}

	// check reactions
	for _, messageReaction := range message.Reactions {
		if messageReaction.Emoji.APIName() == reaction.Emoji.APIName() {
			// lock mutex to prevent quick starring posting a message twice
			// do this as late as possible so we don't lock it unneccessarily
			sb.mu.Lock()
			defer sb.mu.Unlock()
			if messageReaction.Count >= gs.Starboard.ReactLimit {
				err = sb.createOrEditMessage(s, message, reaction.GuildID, messageReaction.Count, messageReaction.Emoji)
				if err != nil {
					sb.Sugar.Errorf("Error creating/editing starboard message for %v: %v", message.ID, err)
				}
				return
			} else if messageReaction.Count < gs.Starboard.ReactLimit {
				if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
					err := sb.deleteStarboardMessage(s, m, reaction.GuildID)
					if err != nil {
						sb.Sugar.Errorf("Error deleting message: %v", err)
					}
					return
				}
			}
		}
	}

	// there were no star reacts, so delete the message
	if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
		err := sb.deleteStarboardMessage(s, m, reaction.GuildID)
		if err != nil {
			sb.Sugar.Errorf("Error deleting message: %v", err)
		}
		return
	}
}

// ReactionRemove ...
func (sb *Sb) ReactionRemove(s *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	// if in DMs, do nothing
	if reaction.GuildID == "" {
		return
	}

	gs, err := sb.Pool.GetGuildSettings(reaction.GuildID)
	if err != nil {
		return
	}

	// if the channel is the server's starboard channel, do nothing
	if reaction.ChannelID == gs.Starboard.StarboardID {
		return
	}
	// if the emoji is not the server's starboard emoji, do nothing
	if reaction.Emoji.MessageFormat() != gs.Starboard.Emoji {
		return
	}

	if sb.Pool.InStarboardBlacklist(reaction.GuildID, reaction.ChannelID) {
		return
	}

	// get the message
	message, err := s.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		sb.Sugar.Errorf("Error getting message %v: %v", reaction.MessageID, err)
	}

	if message.Author.ID == reaction.UserID && !gs.Starboard.SenderCanReact {
		return
	}

	// check reactions
	for _, messageReaction := range message.Reactions {
		if messageReaction.Emoji.APIName() == reaction.Emoji.APIName() {
			// lock mutex to prevent quick starring posting a message twice
			// do this as late as possible so we don't lock it unneccessarily
			sb.mu.Lock()
			defer sb.mu.Unlock()
			if messageReaction.Count >= gs.Starboard.ReactLimit {
				err = sb.createOrEditMessage(s, message, reaction.GuildID, messageReaction.Count, messageReaction.Emoji)
				if err != nil {
					sb.Sugar.Errorf("Error creating/editing starboard message for %v: %v", message.ID, err)
				}
				return
			} else if messageReaction.Count < gs.Starboard.ReactLimit {
				if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
					err := sb.deleteStarboardMessage(s, m, reaction.GuildID)
					if err != nil {
						sb.Sugar.Errorf("Error deleting message: %v", err)
					}
					return
				}
			}
		}
	}

	// there were no star reacts, so delete the message
	if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
		err := sb.deleteStarboardMessage(s, m, reaction.GuildID)
		if err != nil {
			sb.Sugar.Errorf("Error deleting message: %v", err)
		}
		return
	}
}

// MessageDelete ...
func (sb *Sb) MessageDelete(s *discordgo.Session, message *discordgo.MessageDelete) {
	// lock mutex to prevent quick starring posting a message twice
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if m := sb.Pool.GetStarboardEntry(message.ID); m != "" {
		err := sb.Pool.DeleteStarboardEntry(m)
		if err != nil {
			sb.Sugar.Errorf("Error deleting message entry for %v: %v", message.ID, err)
			return
		}
		log.Printf("Deleted starboard database entry for %v", message.ID)
	} else if m := sb.Pool.GetOrigStarboardMessage(message.ID); m != "" {
		err := sb.Pool.DeleteStarboardEntry(m)
		if err != nil {
			sb.Sugar.Errorf("Error deleting message entry for %v: %v", message.ID, err)
			return
		}
		sb.Sugar.Infof("Deleted starboard database entry for %v", message.ID)
	}
}
