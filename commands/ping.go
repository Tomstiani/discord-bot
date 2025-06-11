package commands

import "github.com/bwmarrin/discordgo"

// RegisterPing registers the ping slash command
func RegisterPing() {
	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Check if the bot is responding",
		},
		Handler: handlePing,
	})
}

// handlePing handles the ping slash command
func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🏓 Pong!",
		},
	})
}
