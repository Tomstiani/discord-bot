package bot

import (
	"log"

	"discord-bot/lfg"

	"github.com/bwmarrin/discordgo"
)

// LFG manager instance
var lfgManager *lfg.Manager

// initLFG initializes the LFG manager
func (b *Bot) initLFG() {
	lfgManager = lfg.New(b.Config)
}

// voiceStateUpdate handles voice state changes (join/leave voice channels)
func (b *Bot) voiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Initialize LFG manager if not already done
	if lfgManager == nil {
		b.initLFG()
	}

	// Ignore bot's own voice state changes
	if vs.UserID == s.State.User.ID {
		return
	}

	// Check if user joined a voice channel (vs.ChannelID != "" means they're in a channel)
	if vs.ChannelID != "" && (vs.BeforeUpdate == nil || vs.BeforeUpdate.ChannelID != vs.ChannelID) {
		b.handleUserJoinedVoice(s, vs)
	}
}

// handleUserJoinedVoice processes when a user joins a voice channel
func (b *Bot) handleUserJoinedVoice(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Get channel information
	channel, err := s.Channel(vs.ChannelID)
	if err != nil {
		log.Printf("Error getting channel info: %v", err)
		return
	}

	// Get user information
	user, err := s.User(vs.UserID)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return
	}

	// Check if this is an LFG channel and handle it
	if lfgManager.IsLFGChannel(channel) {
		lfgManager.HandleUserJoinedLFG(s, user, channel)
	}
}
