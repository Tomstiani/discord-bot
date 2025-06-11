package bot

import (
	"fmt"
	"log"

	"discord-bot/commands"

	"github.com/bwmarrin/discordgo"
)

// interactionCreate handles slash command interactions
func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only handle slash commands
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name

	// Look up the command
	command, exists := commands.Get(commandName)
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown command!",
				Flags:   discordgo.MessageFlagsEphemeral, // Only visible to the user
			},
		})
		return
	}

	// Execute the command
	command.Handler(s, i)
}

// registerSlashCommands registers all slash commands with Discord
func (b *Bot) registerSlashCommands() error {
	// First, clean up any existing commands to avoid duplicates
	err := b.cleanupOldCommands()
	if err != nil {
		log.Printf("Warning: Could not clean up old commands: %v", err)
	}

	definitions := commands.GetDefinitions()

	guildID := b.Config.GuildID
	if guildID != "" {
		fmt.Printf("Registering commands to specific server: %s\n", guildID)
	} else {
		fmt.Println("Registering commands globally (may take up to 1 hour to appear)")
	}

	for _, definition := range definitions {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, guildID, definition)
		if err != nil {
			return fmt.Errorf("cannot create slash command %s: %v", definition.Name, err)
		}
		if guildID != "" {
			fmt.Printf("Registered slash command: /%s (server-specific)\n", definition.Name)
		} else {
			fmt.Printf("Registered slash command: /%s (global)\n", definition.Name)
		}
	}

	return nil
}

// cleanupOldCommands removes existing commands to prevent duplicates
func (b *Bot) cleanupOldCommands() error {
	guildID := b.Config.GuildID

	// Get existing commands
	existingCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	// Delete existing commands
	for _, cmd := range existingCommands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, guildID, cmd.ID)
		if err != nil {
			log.Printf("Could not delete command %s: %v", cmd.Name, err)
		} else {
			fmt.Printf("Cleaned up old command: /%s\n", cmd.Name)
		}
	}

	// Also clean up global commands if we're registering server-specific ones
	if guildID != "" {
		globalCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
		if err == nil {
			for _, cmd := range globalCommands {
				err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID)
				if err != nil {
					log.Printf("Could not delete global command %s: %v", cmd.Name, err)
				} else {
					fmt.Printf("Cleaned up old global command: /%s\n", cmd.Name)
				}
			}
		}
	}

	return nil
}
