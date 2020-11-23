package modutilcommands

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

// Members ...
func Members(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	role, err := ctx.ParseRole(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	members, err := ctx.Session.GuildMembers(ctx.Message.GuildID, "", 1000)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var rMembers roleMembers
	for _, m := range members {
		for _, r := range m.Roles {
			if r == role.ID {
				rMembers = append(rMembers, m)
				break
			}
		}
	}

	sort.Sort(rMembers)

	var mString string

	for i, m := range rMembers {
		nick := m.Nick
		if nick == "" {
			nick = m.User.Username
		}
		mString += fmt.Sprintf("%v. %v\n    %v / %v\n\n", i+1, m.User.ID, m.User.String(), nick)
	}

	if len(mString) <= 2000 {
		if len(mString) == 0 {
			mString = "None"
		}
		_, err = ctx.Embedf("Members of "+role.Name, "```%v```", mString)
		return
	}

	reader := bytes.NewReader([]byte(mString))

	file := discordgo.File{
		Name:   fmt.Sprintf("members-%v-%v.txt", ctx.Message.GuildID, role.ID),
		Reader: reader,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Here you go!",
		Files:   []*discordgo.File{&file},
	})
	return
}

type roleMembers []*discordgo.Member

func (r roleMembers) Len() int      { return len(r) }
func (r roleMembers) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r roleMembers) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{r[i].User.String(), r[j].User.String()})
}
