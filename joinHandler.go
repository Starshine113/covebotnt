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

	guildConf, err := pool.GetGuildSettings(m.GuildID)
	if err != nil {
		sugar.Errorf("Error getting guild settings: %v", err)
	}
	if guildConf.Gatekeeper.GatekeeperChannel != "" {
		var msgB bytes.Buffer
		tmpl, err := template.New("gatekeeper").Parse(guildConf.Gatekeeper.GatekeeperMessage)
		if err != nil {
			sugar.Errorf("Error loading template: %v", err)
			return
		}
		if err := tmpl.Execute(&msgB, m); err != nil {
			sugar.Errorf("Error executing template: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(guildConf.Gatekeeper.GatekeeperChannel, msgB.String())
		if err != nil {
			sugar.Errorf("Error sending message: %v", err)
			return
		}
	}
}
