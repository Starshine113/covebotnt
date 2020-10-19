package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var (
	config           structs.BotConfig
	globalSettings   map[string]guildSettings
	channelBlacklist map[string][]string
	db               *pgxpool.Pool

	messageIDMap, starboardMsgIDMap map[string]string
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
func initialise(token, databaseURL string) (err error) {
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

	if token != "" {
		config.Auth.Token = token
	}
	if databaseURL != "" {
		config.Auth.DatabaseURL = databaseURL
	}
	if os.Getenv("CB_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("CB_DB_URL")
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

	// get starboard messages
	messageIDMap, starboardMsgIDMap, err = getStarboardMessages()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	token := flag.String("token", "", "Override the token in config.toml")
	databaseURL := flag.String("db", "", "Override the database URL in config.toml")
	flag.Parse()

	// run all init routines
	err := initialise(*token, *databaseURL)
	if err != nil {
		panic(err)
	}

	// add message create handler for commands
	dg.AddHandler(messageCreateCommand)

	// add guild create handler to initialise data
	dg.AddHandler(guildJoin)

	// add reaction add/remove handler for starboard
	dg.AddHandler(starboardReactionAdd)
	dg.AddHandler(starboardReactionRemove)
	// add message delete handler for starboard
	dg.AddHandler(starboardMessageDelete)

	// add ready handler
	dg.AddHandler(onReady)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildEmojis | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMembers)

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
