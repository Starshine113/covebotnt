package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const botVersion = "0.23"

var (
	config           structs.BotConfig
	globalSettings   map[string]structs.GuildSettings
	channelBlacklist map[string][]string
	db               *pgxpool.Pool
	pool             *cbdb.Db
	boltDb           *cbdb.BoltDb
	levelCache       *ttlcache.Cache
	gitOut           []byte
	router           *crouter.Router
	startTime        time.Time

	messageIDMap, starboardMsgIDMap map[string]string
	handlerMap                      map[string]func()
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
	pool = &cbdb.Db{Pool: db}
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

	// open Bolt db
	bolt, err := bolt.Open(config.Auth.BoltPath, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	boltDb, err = cbdb.BoltInit(bolt)
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

	// create a cache for XP gain
	levelCache = ttlcache.NewCache()
	levelCache.SkipTTLExtensionOnHit(true)
	levelCache.SetTTL(8 * time.Second)

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

	git := exec.Command("git", "rev-parse", "--short", "HEAD")
	gitOut, _ = git.Output()

	router = crouter.NewRouter()
	// add commands
	addUserCommands()
	addHelperCommands()
	addModCommands()
	addAdminCommands()
	addOwnerCommands()
	// add autoresponses
	addAutoResponses()

	// add message create handler for commands
	dg.AddHandler(messageCreateCommand)

	// add guild create handler to initialise data
	dg.AddHandler(guildJoin)

	// add reaction add/remove handler for starboard
	dg.AddHandler(starboardReactionAdd)
	dg.AddHandler(starboardReactionRemove)
	// add message delete handler for starboard
	dg.AddHandler(starboardMessageDelete)

	// add join handler
	dg.AddHandler(onJoin)

	// add ready handler
	dg.AddHandler(onReady)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildEmojis | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMembers)

	// Get start time for uptime command
	startTime = time.Now()

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
		boltDb.Bolt.Close()
		levelCache.Close()
	}()

	sugar.Infof("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
