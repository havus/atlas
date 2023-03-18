package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var Commands = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"inspect-user": inspectUser,
}

func inspectUser(discordSession *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := interactionCreate.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	userTarget := optionMap["user-target"]
	user := userTarget.UserValue(discordSession)
	fmt.Println("user value", user.Email, user.ID)

	data := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("ID: %s\nUsername: %s\nTag: %s#%s", user.ID, user.Username, user.Username, user.Discriminator),
		},
	}

	err := discordSession.InteractionRespond(interactionCreate.Interaction, &data)
	if err != nil {
		log.Printf("failed to interact response %v", err)
	}
}
