package crouter

import (
	"fmt"
	"time"

	"codeberg.org/eviedelta/dwhook"
	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/Starshine113/covebotnt/structs"
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
	case *discordgo.RESTError:
		e := err.(*discordgo.RESTError)
		if e.Message != nil {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v```", e.Message.Message),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Error code: %v", e.Message.Code),
				},
				Color:     0xbf1122,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		} else {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v```", e.ResponseBody),
				Color:       0xbf1122,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			})
		}
		return nil, err
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

		config := ctx.AdditionalParams["config"].(structs.BotConfig)

		if config.Bot.LogWebhook != "" {
			msg := dwhook.Message{
				Content:   fmt.Sprintf("> An internal error occured in %v (%v) of guild %v\n> Triggered by %v (%v/%v):", ctx.Channel.ID, ctx.Channel.Mention(), ctx.Channel.GuildID, ctx.Author.String(), ctx.Author.Mention(), ctx.Author.ID),
				Username:  ctx.BotUser.Username + " Error",
				AvatarURL: ctx.BotUser.AvatarURL("256"),
				Embeds: []dwhook.Embed{{
					Color:       0xbf1122,
					Description: fmt.Sprintf("```%v```", err.Error()),
					Fields: []dwhook.EmbedField{{
						Name:  "Command",
						Value: fmt.Sprintf("**Command**: `%v`\n**Arguments**: `%v`", ctx.Command, ctx.Args),
					}},
					Footer: dwhook.EmbedFooter{
						Text: "Triggering message ID: " + ctx.Author.ID,
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				}},
			}
			dwhook.SendTo(config.Bot.LogWebhook, msg)
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
