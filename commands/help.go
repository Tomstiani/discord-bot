package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RegisterHelp registers the help slash command
func RegisterHelp() {
	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Show all available commands",
		},
		Handler: handleHelp,
	})
}

// handleHelp handles the help slash command
func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var helpText strings.Builder
	helpText.WriteString("**Available Slash Commands:**\n\n")

	for _, cmd := range GetAll() {
		helpText.WriteString(fmt.Sprintf("**/%s** - %s\n",
			cmd.Definition.Name,
			cmd.Definition.Description))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpText.String(),
		},
	})
}
