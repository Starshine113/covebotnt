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
		DMWebhook    string   `toml:"dm_webhook"`
		BlockedUsers []string `toml:"blocked_users"`
		LogWebhook   string   `toml:"log_webhook"`
		CustomStatus struct {
			Type   string `toml:"type"`
			Status string `toml:"status"`
		} `toml:"custom_status"`
	} `toml:"bot"`
	Logging struct {
		LogLevel string `toml:"log_level"`
		DebugURL string `toml:"debug_url"`
		InfoURL  string `toml:"info_url"`
		WarnURL  string `toml:"warn_url"`
		ErrorURL string `toml:"error_url"`
		PanicURL string `toml:"panic_url"`
	} `toml:"logging"`
}
