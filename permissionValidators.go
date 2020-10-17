package main

func checkOwner(userID string) error {
	for _, ownerID := range config.Bot.BotOwners {
		if userID == ownerID {
			return nil
		}
	}
	return &errorNoPermissions{"BotOwner"}
}
