package ownercommands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Guilds ...
func Guilds(ctx *crouter.Ctx) (err error) {
	guilds := make([]string, 0)
	for _, g := range ctx.Session.State.Guilds {
		guilds = append(guilds, fmt.Sprintf("%v (%v)", g.Name, g.ID))
	}

	if len(strings.Join(guilds, "\n")) < 2000 {
		_, err = ctx.Embedf(fmt.Sprintf("Guilds (%v)", len(ctx.Session.State.Guilds)), "```%v```", strings.Join(guilds, "\n"))
		return
	}

	reader := bytes.NewReader([]byte(strings.Join(guilds, "\n")))

	file := discordgo.File{
		Name:   fmt.Sprintf("guilds-%v.txt", time.Now().UTC().Format("2006-01-02-15-04-05")),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})
	return
}
