package crouter

import (
	"fmt"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

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
		id, cmdErr := uuid.NewRandom()
		if cmdErr != nil {
			return err, cmdErr
		}
		cmdErr = ctx.BoltDb.AddError(cbdb.CmdError{
			ErrorID: id.String(),
			Error:   err.Error(),
		})
		if cmdErr != nil {
			return err, cmdErr
		}
		msgErr := ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, ErrorEmoji)
		if msgErr != nil {
			return err, msgErr
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Internal error occured",
			Description: fmt.Sprintf("```%v```\nIf this error persists, please contact the bot developer with the ID above.", err.Error()),
			Color:       0xbf1122,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: id.String(),
			},
		}

		_, msgErr = ctx.Send(&discordgo.MessageSend{
			Embed:   embed,
			Content: "> Error code: `" + id.String() + "`",
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return err, msgErr
	}
	return nil, nil
}
