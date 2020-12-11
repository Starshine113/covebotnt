package structs

// GuildSettings holds a guild's configuration
type GuildSettings struct {
	Prefixes   []string
	Starboard  StarboardSettings
	Moderation ModSettings
	Gatekeeper GatekeeperSettings
	YAG        YAGImport
}

// YAGImport holds the settings for importing mod logs from YAGPDB
type YAGImport struct {
	Channel string
	Enabled bool
}

// StarboardSettings holds the starboard settings
type StarboardSettings struct {
	StarboardID      string
	ReactLimit       int
	Emoji            string
	SenderCanReact   bool
	ReactToStarboard bool
}

// ModSettings holds the mod settings
type ModSettings struct {
	ModRoles    []string
	HelperRoles []string
	ModLog      string
	MuteRole    string
	PauseRole   string
}

// GatekeeperSettings holds the gatekeeper settings
type GatekeeperSettings struct {
	GatekeeperRoles   []string
	MemberRoles       []string
	GatekeeperChannel string
	GatekeeperMessage string
	WelcomeChannel    string
	WelcomeMessage    string
}
