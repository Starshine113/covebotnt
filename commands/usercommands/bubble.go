package usercommands

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/starshine-sys/flagparser"
)

func bubble(ctx *crouter.Ctx) (err error) {
	flagParser, err := flagparser.NewFlagParser(flagparser.Bool("prepop"))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	args, err := flagParser.Parse(ctx.Args)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	size := 10
	pop := args["prepop"].(bool)

	a := args["rest"].([]string)

	for _, arg := range a {
		if strings.ToLower(arg) != "-prepop" {
			size, err = strconv.Atoi(arg)
			if err != nil {
				size = 10
			}
			break
		}
	}

	if size > 13 {
		_, err = ctx.Sendf("%v A size of %v is too large, maximum 13.", crouter.ErrorEmoji, size)
		return
	} else if size < 1 {
		_, err = ctx.Sendf("%v A size of %v is too small, minimum 1.", crouter.ErrorEmoji, size)
		return
	}

	var out string
	for i := size; i > 0; i-- {
		for j := size; j > 0; j-- {
			if pop {
				if j != 1 && j != size && i != 1 && i != size {
					if rand.Intn(6) == 5 {
						out += "pop"
					} else {
						out += "||pop||"
					}
				} else {
					out += "||pop||"
				}
			} else {
				out += "||pop||"
			}
		}
		out += "\n"
	}
	_, err = ctx.Send(out)
	return err
}
