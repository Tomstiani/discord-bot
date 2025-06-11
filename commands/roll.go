package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	// Seed random number generator when package loads
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RegisterRoll registers the roll slash command
func RegisterRoll() {
	minValue := float64(1)
	maxValue := float64(1000)

	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "roll",
			Description: "Roll a dice with specified number of sides",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "sides",
					Description: "Number of sides on the dice (default: 6)",
					Required:    false,
					MinValue:    &minValue,
					MaxValue:    maxValue,
				},
			},
		},
		Handler: handleRoll,
	})
}

// handleRoll handles the roll slash command
func handleRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sides := 6 // Default to 6-sided die

	// Get the sides parameter if provided
	if len(i.ApplicationCommandData().Options) > 0 {
		for _, option := range i.ApplicationCommandData().Options {
			if option.Name == "sides" {
				sides = int(option.IntValue())
			}
		}
	}

	// Roll the dice
	result := rand.Intn(sides) + 1

	response := fmt.Sprintf("🎲 You rolled a **%d** out of %d!", result, sides)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
