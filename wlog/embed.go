package wlog

import (
	"fmt"

	"codeberg.org/eviedelta/dwhook"
)

func embed(title, err string, colour int) dwhook.Embed {
	return dwhook.Embed{
		Title:       title,
		Description: fmt.Sprintf("```%v```", err),
		Color:       colour,
	}
}

func (w *Wlog) send(url, level string, embeds ...dwhook.Embed) (err error) {
	_, err = dwhook.SendTo(url, dwhook.Message{
		Username:  fmt.Sprintf("%v %v", w.Name, level),
		AvatarURL: w.AvatarURL,
		Embeds:    embeds,
	})
	return err
}
