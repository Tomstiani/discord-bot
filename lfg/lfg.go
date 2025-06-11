package lfg

import (
	"fmt"
	"log"
	"strings"

	"discord-bot/config"

	"github.com/bwmarrin/discordgo"
)

// Manager handles all LFG (Looking for Game) functionality
type Manager struct {
	Config *config.Config
}

// New creates a new LFG manager
func New(cfg *config.Config) *Manager {
	return &Manager{
		Config: cfg,
	}
}

// IsLFGChannel checks if a channel is designated for "Looking for Game"
func (m *Manager) IsLFGChannel(channel *discordgo.Channel) bool {
	lfgChannelId := m.Config.LFGChannelID
	if lfgChannelId != "" && lfgChannelId == channel.ID {
		return true
	}
	return false
}

// HandleUserJoinedLFG processes when someone joins an LFG channel
func (m *Manager) HandleUserJoinedLFG(s *discordgo.Session, user *discordgo.User, channel *discordgo.Channel) {
	fmt.Printf("🎮 %s joined %s - Looking for game!\n", user.Username, channel.Name)

	// Send a message to announce the LFG
	m.announceUserLookingForGame(s, user, channel)

	// TODO: Add game selection interface
	// TODO: Add NTFY notifications to subscribers
}

// announceUserLookingForGame sends a message when someone is looking for a game
func (m *Manager) announceUserLookingForGame(s *discordgo.Session, user *discordgo.User, voiceChannel *discordgo.Channel) {
	message := fmt.Sprintf("@everyone 🎮 **%s** is looking for people to play! What game do you want to play?", user.Username)

	// Use configured announcement channel or find one automatically
	var textChannelID string
	if m.Config.LFGAnnouncementChannelID != "" {
		textChannelID = m.Config.LFGAnnouncementChannelID
		fmt.Printf("📢 Using configured announcement channel: %s\n", textChannelID)
	} else {
		textChannelID = m.findAnnouncementChannel(s, voiceChannel.GuildID)
		fmt.Printf("📢 Auto-found announcement channel: %s\n", textChannelID)
	}

	if textChannelID != "" {
		_, err := s.ChannelMessageSend(textChannelID, message)
		if err != nil {
			log.Printf("Error sending LFG announcement: %v", err)
		} else {
			fmt.Printf("📢 Sent LFG announcement with @everyone tag\n")
		}
	} else {
		log.Println("Warning: Could not find a suitable text channel for LFG announcement")
	}
}

// findAnnouncementChannel finds the best text channel to send LFG announcements
func (m *Manager) findAnnouncementChannel(s *discordgo.Session, guildID string) string {
	guild, err := s.Guild(guildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return ""
	}

	var textChannelID string
	var fallbackChannelID string

	for _, ch := range guild.Channels {
		if ch.Type == discordgo.ChannelTypeGuildText {
			// Store first text channel as fallback
			if fallbackChannelID == "" {
				fallbackChannelID = ch.ID
			}

			// Prefer gaming, general, or LFG channels
			channelNameLower := strings.ToLower(ch.Name)
			if strings.Contains(channelNameLower, "general") ||
				strings.Contains(channelNameLower, "gaming") ||
				strings.Contains(channelNameLower, "lfg") ||
				strings.Contains(channelNameLower, "announcements") {
				textChannelID = ch.ID
				break
			}
		}
	}

	// Use preferred channel if found, otherwise use fallback
	if textChannelID != "" {
		return textChannelID
	}
	return fallbackChannelID
}

// TODO: Future methods to add:
// - SelectGame(user, game) - Let user select what game they want to play
// - NotifySubscribers(game, user) - Send NTFY notifications to game subscribers
// - GetActiveSession(channelID) - Get current LFG session for a channel
// - CreateGameSession(user, game, channel) - Create a new game session
