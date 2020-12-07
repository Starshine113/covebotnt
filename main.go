package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/covebotnt/bot"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/commands/admincommands"
	"github.com/Starshine113/covebotnt/commands/modcommands"
	"github.com/Starshine113/covebotnt/commands/modutilcommands"
	"github.com/Starshine113/covebotnt/commands/ownercommands"
	"github.com/Starshine113/covebotnt/commands/usercommands"
	"github.com/Starshine113/covebotnt/crouter"
	"github.com/Starshine113/covebotnt/notes"
	"github.com/Starshine113/covebotnt/starboard"
	"github.com/Starshine113/covebotnt/structs"
	"github.com/Starshine113/covebotnt/triggers"
	"github.com/Starshine113/covebotnt/wlog"
	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

const botVersion = "0.99.999"

var (
	config    structs.BotConfig
	pool      *cbdb.Db
	boltDb    *cbdb.BoltDb
	gitOut    []byte
	router    *crouter.Router
	startTime time.Time
	b         *bot.Bot

	handlerMap *ttlcache.Cache
)

var (
	sugar *wlog.Wlog
	dg    *discordgo.Session
)

func loadDB(c structs.BotConfig) (err error) {
	pool, err = cbdb.DbInit(config)
	return err
}

// basic initialisation routines
func initialise(token, databaseURL string) (err error) {
	// load config
	config, err = getConfig()
	if err != nil {
		return err
	}

	sugar = wlog.Logger(wlog.URLs{
		DebugURL: config.Logging.DebugURL,
		InfoURL:  config.Logging.InfoURL,
		WarnURL:  config.Logging.WarnURL,
		ErrorURL: config.Logging.ErrorURL,
		PanicURL: config.Logging.PanicURL,
	}, config.Logging.LogLevel)

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
	err = loadDB(config)
	if err != nil {
		return err
	}

	handlerMap = ttlcache.NewCache()
	handlerMap.SetCacheSizeLimit(10000)
	handlerMap.SetTTL(15 * time.Minute)
	handlerMap.SetExpirationCallback(func(key string, value interface{}) {
		value.(func())()
	})

	git := exec.Command("git", "rev-parse", "--short", "HEAD")
	gitOut, _ = git.Output()

	b = bot.NewBot(dg, sugar, pool, boltDb, config, handlerMap, botVersion, string(gitOut), time.Now())

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

	sb := starboard.Sb{
		Sugar: sugar,
		Pool:  pool,
	}

	router = crouter.NewRouter(b)
	// add commands
	modcommands.Init(router)
	modutilcommands.Init(router)
	ownercommands.Init(router)
	usercommands.Init(router)
	admincommands.Init(router)
	notes.Init(router)
	triggers.Init(router)

	addStarboardCommands()
	addGkCommands()

	addAdminCommands()
	addOwnerCommands()

	// add autoresponses
	addAutoResponses()

	// add message create handler for commands
	dg.AddHandler(router.MessageCreate)

	// add guild create handler to initialise data
	dg.AddHandler(guildJoin)

	// add reaction add/remove handler for starboard
	dg.AddHandler(sb.ReactionAdd)
	dg.AddHandler(sb.ReactionRemove)
	// add message delete handler for starboard
	dg.AddHandler(sb.MessageDelete)

	// add join handler
	dg.AddHandler(onJoin)

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
		pool.Pool.Close()
		sugar.Infof("Closed database connection.")

		sugar.Flush()
		boltDb.Bolt.Close()
	}()

	sugar.Infof("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
