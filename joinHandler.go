package main

import (
	"bytes"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

func onJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// if the user is a bot, return
	if m.User.Bot {
		return
	}

	guildConf := globalSettings[m.GuildID]
	if guildConf.Gatekeeper.GatekeeperChannel != "" {
		var msgB bytes.Buffer
		tmpl, err := template.New("gatekeeper").Parse(guildConf.Gatekeeper.GatekeeperMessage)
		if err != nil {
			sugar.Error(err)
			return
		}
		if err := tmpl.Execute(&msgB, m); err != nil {
			sugar.Error(err)
			return
		}

		_, err = s.ChannelMessageSend(guildConf.Gatekeeper.GatekeeperChannel, msgB.String())
		if err != nil {
			sugar.Error(err)
			return
		}
	}
}
