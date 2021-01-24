package modutilcommands

import (
	"fmt"

	"github.com/starshine-sys/covebotnt/crouter"
)

// RefreshMVC refreshes the mvc role
func RefreshMVC(ctx *crouter.Ctx) (err error) {
	mvc, err := ctx.ParseRole("mvc")
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	mvcOld, err := ctx.ParseRole("mvc-old")
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	members, err := ctx.Session.GuildMembers(ctx.Message.GuildID, "", 1000)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var count int64 = 0
	var mvcUsers []string
	for _, m := range members {
		for _, r := range m.Roles {
			if r == mvc.ID {
				count++
				mvcUsers = append(mvcUsers, m.User.ID)
				break
			}
		}
	}

	for _, m := range mvcUsers {
		err = ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, m, mvcOld.ID)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
	}

	for _, m := range mvcUsers {
		err = ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, m, mvc.ID)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
	}

	_, err = ctx.Embed("Success", fmt.Sprintf("Refreshed the mvc role. Previous mvc count: %v.", count), 0x21a1a8)
	return
}
