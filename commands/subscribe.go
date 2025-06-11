package commands

import (
	"fmt"
	"strings"

	"discord-bot/data"

	"github.com/bwmarrin/discordgo"
)

// Common games list - you can expand this
var commonGames = []string{
	"valorant", "csgo", "cs2", "overwatch", "apex", "fortnite",
	"minecraft", "rocket-league", "cod", "warzone", "dota2", "lol",
	"among-us", "fall-guys", "gta", "rust", "destiny2", "wow",
}

// Global subscription manager - you'll initialize this in main
var SubManager *data.SubscriptionManager

// Simple title case function
func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// RegisterSubscribe registers the subscribe slash command
func RegisterSubscribe() {
	// Create choices for common games
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(commonGames))
	for _, game := range commonGames {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  toTitleCase(strings.ReplaceAll(game, "-", " ")),
			Value: game,
		})
	}

	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "subscribe",
			Description: "Subscribe to notifications for a game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "game",
					Description: "The game you want notifications for",
					Required:    true,
					Choices:     choices,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ntfy-topic",
					Description: "Your NTFY topic (e.g., 'john_gaming')",
					Required:    true,
				},
			},
		},
		Handler: handleSubscribe,
	})
}

// RegisterUnsubscribe registers the unsubscribe slash command
func RegisterUnsubscribe() {
	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "unsubscribe",
			Description: "Unsubscribe from game notifications",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "game",
					Description:  "The game to unsubscribe from",
					Required:     true,
					Autocomplete: true, // We'll make this show user's subscriptions
				},
			},
		},
		Handler: handleUnsubscribe,
	})
}

// RegisterMyGames registers the mygames slash command
func RegisterMyGames() {
	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "mygames",
			Description: "See what games you're subscribed to",
		},
		Handler: handleMyGames,
	})
}

// RegisterGamesList registers the games slash command
func RegisterGamesList() {
	Register(&SlashCommand{
		Definition: &discordgo.ApplicationCommand{
			Name:        "games",
			Description: "See all games people are subscribed to",
		},
		Handler: handleGamesList,
	})
}

// handleSubscribe handles the subscribe command
func handleSubscribe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var game, ntfyTopic string

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "game":
			game = option.StringValue()
		case "ntfy-topic":
			ntfyTopic = option.StringValue()
		}
	}

	user := i.Member.User
	err := SubManager.Subscribe(user.ID, user.Username, game, ntfyTopic)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Error: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	gameName := toTitleCase(strings.ReplaceAll(game, "-", " "))
	response := fmt.Sprintf("✅ Successfully subscribed to **%s** notifications!\n"+
		"NTFY Topic: `%s`\n"+
		"You'll get notified when someone wants to play!", gameName, ntfyTopic)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleUnsubscribe handles the unsubscribe command
func handleUnsubscribe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	game := i.ApplicationCommandData().Options[0].StringValue()
	user := i.Member.User

	err := SubManager.Unsubscribe(user.ID, game)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Error: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	gameName := toTitleCase(strings.ReplaceAll(game, "-", " "))
	response := fmt.Sprintf("✅ Successfully unsubscribed from **%s** notifications!", gameName)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleMyGames handles the mygames command
func handleMyGames(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	subscriptions := SubManager.GetSubscriptions(user.ID)

	if len(subscriptions) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "📱 You're not subscribed to any games yet!\nUse `/subscribe` to get started.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var response strings.Builder
	response.WriteString("📱 **Your Game Subscriptions:**\n\n")
	for _, sub := range subscriptions {
		gameName := toTitleCase(strings.ReplaceAll(sub.Game, "-", " "))
		response.WriteString(fmt.Sprintf("🎮 **%s** → `%s`\n", gameName, sub.NTFYTopic))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleGamesList handles the games command
func handleGamesList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	games := SubManager.GetAllGames()

	if len(games) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🎮 No one is subscribed to any games yet!",
			},
		})
		return
	}

	var response strings.Builder
	response.WriteString("🎮 **Games with subscribers:**\n\n")
	for _, game := range games {
		subscribers := SubManager.GetSubscribersForGame(game)
		gameName := toTitleCase(strings.ReplaceAll(game, "-", " "))
		response.WriteString(fmt.Sprintf("**%s** (%d subscribers)\n", gameName, len(subscribers)))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response.String(),
		},
	})
}
