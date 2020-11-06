package cbctx

import "github.com/bwmarrin/discordgo"

// CommandError sends an error message and optionally returns an error for logging purposes
func (ctx *Ctx) CommandError(err error) (error, error) {
	switch err.(type) {
	case *ErrorNoPermissions, *ErrorNoDMs:
		ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " You are not allowed to use this command:\n> " + err.Error(),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		if msgErr != nil {
			return nil, msgErr
		}
	case *ErrorMissingRequiredArgs, *ErrorNotEnoughArgs:
		ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " Command call was missing arguments:\n> " + err.Error(),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		if msgErr != nil {
			return nil, msgErr
		}
	default:
		ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, ErrorEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: ErrorEmoji + " An internal error occured:\n> " + err.Error(),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return err, msgErr
	}
	return nil, nil
}
