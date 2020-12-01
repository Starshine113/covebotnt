package wlog

import (
	"fmt"

	"codeberg.org/eviedelta/dwhook"
)

// Errorf ...
func (w *Wlog) Errorf(template string, args ...interface{}) {
	if w.logLevel > logLevelError {
		return
	}

	msg := template
	if len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	if w.urls.ErrorURL != "" {
		w.send(w.urls.ErrorURL, "Error", embed("Error", msg, 0xea5a2e))
	}

	w.sugar.Errorf(template, args...)
}

// Infof ...
func (w *Wlog) Infof(template string, args ...interface{}) {
	if w.logLevel > logLevelInfo {
		return
	}

	msg := template
	if len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	if w.urls.InfoURL != "" {
		w.infoMu.Lock()
		w.infoBuffer = append(w.infoBuffer, embed("Info", msg, 0x93ea2e))

		if len(w.infoBuffer) > 2 {
			w.send(w.urls.InfoURL, "Info", w.infoBuffer...)
			w.infoBuffer = make([]dwhook.Embed, 0)
		}
		w.infoMu.Unlock()
	}

	w.sugar.Infof(template, args...)
}

// Debugf ...
func (w *Wlog) Debugf(template string, args ...interface{}) {
	if w.logLevel > logLevelDebug {
		return
	}

	msg := template
	if len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	if w.urls.DebugURL != "" {
		w.debugMu.Lock()
		w.debugBuffer = append(w.debugBuffer, embed("Debug", msg, 0x93ea2e))

		if len(w.infoBuffer) > 2 {
			w.send(w.urls.DebugURL, "Debug", w.debugBuffer...)
			w.debugBuffer = make([]dwhook.Embed, 0)
		}
		w.debugMu.Unlock()
	}

	w.sugar.Debugf(template, args...)
}

// Panicf ...
func (w *Wlog) Panicf(template string, args ...interface{}) {
	msg := template
	if len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	if w.urls.PanicURL != "" {
		w.send(w.urls.PanicURL, "Panic", embed("Panic", msg, 0xea5a2e))
	}

	w.sugar.Panicf(template, args...)
}
