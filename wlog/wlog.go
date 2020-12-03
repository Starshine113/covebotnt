package wlog

import (
	"sync"
	"time"

	"codeberg.org/eviedelta/dwhook"
	"go.uber.org/zap"
)

const (
	bufferSize = 3
)

// Wlog is a wrapper for Zap's sugared logger that also logs to Discord webhooks
type Wlog struct {
	sugar       *zap.SugaredLogger
	logLevel    logLevel
	urls        *URLs
	Name        string
	AvatarURL   string
	debugMu     sync.Mutex
	debugBuffer []dwhook.Embed
	infoMu      sync.Mutex
	infoBuffer  []dwhook.Embed
}

// URLs ...
type URLs struct {
	DebugURL string
	InfoURL  string
	WarnURL  string
	ErrorURL string
	PanicURL string
	LogLevel string
}

type logLevel int

const (
	logLevelDebug logLevel = iota
	logLevelInfo
	logLevelWarn
	logLevelError
	logLevelPanic
)

// Logger creates a new Wlog object
func Logger(c URLs, level string) *Wlog {
	wlog := &Wlog{}
	log, _ := zap.NewDevelopment()

	wlog.debugBuffer = make([]dwhook.Embed, 0)
	wlog.infoBuffer = make([]dwhook.Embed, 0)

	wlog.sugar = log.Sugar()
	wlog.urls = &c
	wlog.logLevel = stringToLevel(level)
	wlog.Name = "Bot"
	wlog.AvatarURL = "https://cdn.discordapp.com/embed/avatars/0.png"

	go func() {
		for {
			time.Sleep(30 * time.Second)
			wlog.Flush()
		}
	}()

	return wlog
}

func stringToLevel(s string) logLevel {
	switch s {
	case "PANIC", "panic":
		return logLevelPanic
	case "ERROR", "error", "err", "ERR":
		return logLevelError
	case "WARN", "warn":
		return logLevelWarn
	case "INFO", "info":
		return logLevelInfo
	case "DEBUG", "debug":
	default:
		return logLevelDebug
	}

	return logLevelDebug
}
