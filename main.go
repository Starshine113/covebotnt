package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	config            = getConfig()
	db                = initDB()
	globalSettings, _ = getSettingsAll()
	channelBlacklist  = getBlacklistAll()
)

var (
	logger, _ = zap.NewDevelopment()
	sugar     = logger.Sugar()
	dg, _     = discordgo.New("Bot " + config.Auth.Token)
)

func main() {
	defer logger.Sync()

	// add message create handler for commands
	dg.AddHandler(messageCreateCommand)

	// set intents
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildEmojis | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessageReactions)

	err := dg.Open()
	if err != nil {
		panic(err)
	}

	sugar.Infof("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")
	dg.UpdateStatus(0, config.Bot.Prefixes[0]+"help")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
	dg.Close()
	sugar.Infof("Disconnected from Discord.")
	db.Close()
	sugar.Infof("Closed database connection.")

	os.Exit(0)
}
