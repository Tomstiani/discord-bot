package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the bot
type Config struct {
	Token                    string
	GuildID                  string
	LFGChannelID             string
	LFGAnnouncementChannelID string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get the bot token from environment variable
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("No token provided. Set DISCORD_BOT_TOKEN environment variable.")
	}

	// Get the guild ID from environment variable
	guildID := os.Getenv("DISCORD_GUILD_ID")
	if guildID == "" {
		log.Println("Warning: No DISCORD_GUILD_ID provided. Commands will be registered globally (takes up to 1 hour to appear).")
	}

	// Get the LFG channel ID from environment variable
	lfgChannelID := os.Getenv("DISCORD_LFG_CHANNEL_ID")
	if lfgChannelID == "" {
		log.Println("Warning: No DISCORD_LFG_CHANNEL_ID provided. LFG monitoring will be disabled.")
	}

	// Get the LFG announcement channel ID from environment variable
	lfgAnnouncementChannelID := os.Getenv("DISCORD_LFG_ANNOUNCEMENT_CHANNEL_ID")
	if lfgAnnouncementChannelID == "" {
		log.Println("Warning: No DISCORD_LFG_ANNOUNCEMENT_CHANNEL_ID provided. Will try to find a suitable channel automatically.")
	}

	return &Config{
		Token:                    token,
		GuildID:                  guildID,
		LFGChannelID:             lfgChannelID,
		LFGAnnouncementChannelID: lfgAnnouncementChannelID,
	}
}
