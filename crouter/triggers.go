package crouter

import "strings"

func (r *Router) triggers(channelID, guildID, content string) {
	triggers, err := r.Bot.Pool.Triggers(guildID)
	if err != nil {
		return
	}
	if len(triggers) == 0 {
		return
	}

	for _, t := range triggers {
		if strings.ToLower(content) == strings.ToLower(t.Trigger) {
			r.Bot.Session.ChannelMessageSend(channelID, t.Response)
		}
	}
}
