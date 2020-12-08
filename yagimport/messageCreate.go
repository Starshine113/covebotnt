package yagimport

import (
	"regexp"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
)

// just hardcode the ID it shouldn't change
const yagID = "204255221017214977"

// regex for getting the moderator
var modRegexp = regexp.MustCompile(`\(ID (\d{16,})\)`)

// regex for getting the duration
var durationRegexp = regexp.MustCompile(`Duration: ([\w\s]+)`)

// regexes for figuring out what the action type is + capturing ID and reason
var warnRegexp = regexp.MustCompile(`\*\*⚠Warned [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*)`)
var muteRegexp = regexp.MustCompile(`\*\*🔇Muted [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*) \(\[`)
var unmuteRegexp = regexp.MustCompile(`\*\*🔊Unmuted [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*)`)

func (y *yag) messageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	// only messages in guilds matter
	if m.GuildID == "" {
		return
	}

	// get guild settings
	gs, err := y.Pool.GetGuildSettings(m.GuildID)
	if err != nil {
		return
	}

	if m.ChannelID != gs.YAG.Channel || !gs.YAG.Enabled {
		return
	}

	y.process(m.GuildID, m, gs)
}

func (y *yag) process(g string, m *discordgo.MessageCreate, gs structs.GuildSettings) {
	// no embeds so we don't care
	if len(m.Embeds) == 0 {
		return
	}

	m.GuildID = g

	entry := y.parseEmbed(m.GuildID, m.Embeds[0])
	if entry == nil {
		return
	}

	ts, err := discordgo.SnowflakeTimestamp(m.ID)
	entry.Time = ts

	entry, err = y.Pool.AddToModLog(entry)
	if err != nil {
		y.Sugar.Errorf("Error adding mod log entry: %v", err)
	} else {
		y.Sugar.Debugf("Added mod log entry %v", entry.ID)
	}
}

func (y *yag) parseEmbed(g string, e *discordgo.MessageEmbed) *cbdb.ModLogEntry {
	if warnRegexp.MatchString(e.Description) {
		return y.warnEmbed(g, e)
	}

	if muteRegexp.MatchString(e.Description) {
		return y.muteEmbed(g, e)
	}

	if unmuteRegexp.MatchString(e.Description) {
		return y.unmuteEmbed(g, e)
	}

	return nil
}

func (y *yag) warnEmbed(g string, e *discordgo.MessageEmbed) *cbdb.ModLogEntry {
	groups := warnRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user := groups[1]
	reason := groups[2]

	groups = modRegexp.FindStringSubmatch(e.Author.Name)
	mod := groups[1]

	return &cbdb.ModLogEntry{
		GuildID: g,
		UserID:  user,
		ModID:   mod,
		Type:    "warn",
		Reason:  reason,
	}
}

func (y *yag) muteEmbed(g string, e *discordgo.MessageEmbed) *cbdb.ModLogEntry {
	groups := muteRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user := groups[1]
	reason := groups[2]

	groups = modRegexp.FindStringSubmatch(e.Author.Name)
	mod := groups[1]

	if e.Footer != nil {
		if durationRegexp.MatchString(e.Footer.Text) {
			groups = durationRegexp.FindStringSubmatch(e.Footer.Text)
			if groups[1] != "" && groups[1] != "permanent" {
				reason += " (duration: " + groups[1] + ")"
			}
		}
	}

	return &cbdb.ModLogEntry{
		GuildID: g,
		UserID:  user,
		ModID:   mod,
		Type:    "mute",
		Reason:  reason,
	}
}

func (y *yag) unmuteEmbed(g string, e *discordgo.MessageEmbed) *cbdb.ModLogEntry {
	groups := unmuteRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user := groups[1]
	reason := groups[2]

	groups = modRegexp.FindStringSubmatch(e.Author.Name)
	mod := groups[1]

	return &cbdb.ModLogEntry{
		GuildID: g,
		UserID:  user,
		ModID:   mod,
		Type:    "unmute",
		Reason:  reason,
	}
}
