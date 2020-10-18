package main

import "github.com/bwmarrin/discordgo"

func commandError(err error, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.MessageReactionAdd(m.ChannelID, m.ID, errorEmoji)

	switch err.(type) {
	case *errorNoPermissions, *errorNoDMs:
		_, msgErr := s.ChannelMessageSend(m.ChannelID, errorEmoji+" You are not allowed to use this command:\n> "+err.Error())
		if msgErr != nil {
			sugar.Errorf("An error occured while sending the error message", msgErr)
		}
	case *errorMissingRequiredArgs, *errorNotEnoughArgs:
		_, msgErr := s.ChannelMessageSend(m.ChannelID, errorEmoji+" Command call was missing arguments:\n> "+err.Error())
		if msgErr != nil {
			sugar.Errorf("An error occured while sending the error message", msgErr)
		}
	default:
		sugar.Errorf("Command error occured: ", err.Error())
		_, msgErr := s.ChannelMessageSend(m.ChannelID, errorEmoji+" An internal error occured:\n> "+err.Error())
		if msgErr != nil {
			sugar.Errorf("An error occured while sending the error message", msgErr)
		}
	}
}
