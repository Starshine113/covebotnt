package commands

import (
	"math"
	"strconv"

	"github.com/Starshine113/covebotnt/crouter"
)

// MajorityOf gives the majority of the given number, with optional abstains
func MajorityOf(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckArgRange(1, 2); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var count, abstains int64 = 0, 0
	count, err = strconv.ParseInt(ctx.Args[0], 0, 0)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	if len(ctx.Args) == 2 {
		abstains, err = strconv.ParseInt(ctx.Args[1], 0, 0)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
	}

	majority := (count-abstains)/2 + 1
	_, err = ctx.Embedf("Majority", "The majority of %v with %v abstains is %v.", count, abstains, majority)
	return
}

// MVC gives the majority of the mvc role
func MVC(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckArgRange(0, 1); err != nil {
		_, err = ctx.CommandError(err)
		return err
	}

	var abstains int64 = 0
	if len(ctx.Args) == 1 {
		abstains, err = strconv.ParseInt(ctx.Args[0], 0, 0)
		if err != nil {
			_, err = ctx.CommandError(err)
			return err
		}
	}

	mvc, err := ctx.ParseRole("mvc")
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
	for _, m := range members {
		for _, r := range m.Roles {
			if r == mvc.ID {
				count++
				break
			}
		}
	}

	majority := math.Ceil(float64(((count - abstains) / 2) + 1))
	_, err = ctx.Embedf("MVC majority", "The majority of %v with %v abstains is %v.", count, abstains, majority)
	return
}
