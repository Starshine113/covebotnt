package wlog

import (
	"codeberg.org/eviedelta/dwhook"
)

// Flush sends all buffered embeds
func (w *Wlog) Flush() {
	if w.urls.InfoURL != "" {
		w.infoMu.Lock()
		if len(w.infoBuffer) > 0 {
			w.send(w.urls.InfoURL, "Info", w.infoBuffer...)
			w.infoBuffer = make([]dwhook.Embed, 0)
		}
		w.infoMu.Unlock()
	}
	if w.urls.DebugURL != "" {
		w.debugMu.Lock()
		if len(w.infoBuffer) > 0 {
			w.send(w.urls.DebugURL, "Debug", w.debugBuffer...)
			w.debugBuffer = make([]dwhook.Embed, 0)
		}
		w.debugMu.Unlock()
	}
}
