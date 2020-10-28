package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/Starshine113/covebotnt/cbctx"
)

func commandUpdate(ctx *cbctx.Ctx) (err error) {
	_, err = ctx.Send("Updating Git repository...")
	if err != nil {
		return err
	}

	git := exec.Command("git", "pull")
	pullOutput, err := git.Output()
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	_, err = ctx.Send(fmt.Sprintf("Git:\n```%v```", string(pullOutput)))
	if err != nil {
		return err
	}

	update := exec.Command("/usr/local/go/bin/go", "build")

	outPipe, err := update.StdoutPipe()
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = update.Start()
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	err = update.Wait()
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	updateOutput, err := ioutil.ReadAll(outPipe)
	if err != nil {
		return err
	}
	if len(updateOutput) == 0 {
		updateOutput = []byte("[No command output]")
	}

	_, err = ctx.Send(fmt.Sprintf("`go build`:\n```%v```", string(updateOutput)))
	return
}
