package bot

import (
	"fmt"

	"discord-bot/commands"
	"discord-bot/config"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the Discord bot
type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
}

// New creates a new bot instance
func New(cfg *config.Config) (*Bot, error) {
	// Create a new Discord session
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	bot := &Bot{
		Session: dg,
		Config:  cfg,
	}

	// Register event handlers
	dg.AddHandler(bot.ready)
	dg.AddHandler(bot.interactionCreate)
	dg.AddHandler(bot.voiceStateUpdate)

	// Set required intents for slash commands and voice state updates
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	// Initialize commands
	commands.Initialize()

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	// Register slash commands with Discord
	err = b.registerSlashCommands()
	if err != nil {
		return fmt.Errorf("error registering slash commands: %v", err)
	}

	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.Session.Close()
}

// ready handles the ready event
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Logged in as: %v#%v\n", event.User.Username, event.User.Discriminator)
	fmt.Println("Slash commands registered! Try typing / in Discord to see them.")
}
