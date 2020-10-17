package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var (
	config           botConfig
	globalSettings   map[string]guildSettings
	channelBlacklist map[string][]string
	db               *pgxpool.Pool
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	dg     *discordgo.Session
)

func loadDB() (err error) {
	db, err = initDB()
	if err != nil {
		return err
	}
	channelBlacklist = getBlacklistAll()
	globalSettings, err = getSettingsAll()

	return err
}

// basic initialisation routines
func initialise() (err error) {
	// initialise logger
	logger, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	sugar = logger.Sugar()

	// load config
	config, err = getConfig()
	if err != nil {
		return err
	}

	// get discord session
	dg, err = discordgo.New("Bot " + config.Auth.Token)
	if err != nil {
		return err
	}

	// open database connection and load state
	err = loadDB()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// run all init routines
	err := initialise()
	if err != nil {
		panic(err)
	}

	// add message create handler for commands
	dg.AddHandler(messageCreateCommand)

	// add guild create handler to initialise data
	dg.AddHandler(guildCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildEmojis | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessageReactions)

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		dg.Close()
		sugar.Infof("Disconnected from Discord.")
		db.Close()
		sugar.Infof("Closed database connection.")

		logger.Sync()
	}()

	sugar.Infof("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")
	err = dg.UpdateStatus(0, "testing, use "+config.Bot.Prefixes[0]+"help")
	if err != nil {
		sugar.Errorw("UpdateStatus Error", "Error", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
