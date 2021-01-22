package watchlist

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (uh *uh) onJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// if the user is a bot, return
	if m.User.Bot {
		return
	}

	guildConf, err := uh.pool.GetGuildSettings(m.GuildID)
	if err != nil {
		uh.sugar.Errorf("Error getting guild settings: %v", err)
	}

	if !uh.pool.OnWatchlist(m.GuildID, m.User.ID) {
		return
	}

	if guildConf.Gatekeeper.WatchlistChannel != "" {
		_, err = s.ChannelMessageSendEmbed(guildConf.Gatekeeper.WatchlistChannel, &discordgo.MessageEmbed{
			Title:       "⚠ Person on watchlist joined",
			Description: fmt.Sprintf("%v (%v/%v) just joined the server and is on the watchlist.", m.User.Mention(), m.User.String(), m.User.ID),
			Color:       0xeaa402,
		})
		if err != nil {
			uh.sugar.Errorf("Error sending message: %v", err)
			return
		}
	}
}
