package cbctx

// CommandError sends an error message and optionally returns an error for logging purposes
func (ctx *Ctx) CommandError(err error) (error, error) {
	ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, ErrorEmoji)

	switch err.(type) {
	case *ErrorNoPermissions, *ErrorNoDMs:
		_, msgErr := ctx.Send(ErrorEmoji + " You are not allowed to use this command:\n> " + err.Error())
		if msgErr != nil {
			return nil, msgErr
		}
	case *ErrorMissingRequiredArgs, *ErrorNotEnoughArgs:
		_, msgErr := ctx.Send(ErrorEmoji + " Command call was missing arguments:\n> " + err.Error())
		if msgErr != nil {
			return nil, msgErr
		}
	default:
		_, msgErr := ctx.Send(ErrorEmoji + " An internal error occured:\n> " + err.Error())
		return err, msgErr
	}
	return nil, nil
}
