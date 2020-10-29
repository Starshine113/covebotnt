package structs

// BotConfig holds the bot's configuration
type BotConfig struct {
	Auth struct {
		Token       string `toml:"token"`
		DatabaseURL string `toml:"database_url"`
		BoltPath    string `toml:"bolt_path"`
	} `toml:"auth"`
	Bot struct {
		Prefixes     []string `toml:"prefixes"`
		BotOwners    []string `toml:"bot_owners"`
		Invite       string   `toml:"invite"`
		AllowedBots  []string `toml:"allowed_bots"`
		CustomStatus struct {
			Type   string `toml:"type"`
			Status string `toml:"status"`
		} `toml:"custom_status"`
	} `toml:"bot"`
}
