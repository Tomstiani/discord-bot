package commands

import "github.com/bwmarrin/discordgo"

// SlashCommand represents a slash command
type SlashCommand struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Registry holds all registered slash commands
var Registry = make(map[string]*SlashCommand)

// Register adds a slash command to the registry
func Register(cmd *SlashCommand) {
	Registry[cmd.Definition.Name] = cmd
}

// GetAll returns all registered commands
func GetAll() map[string]*SlashCommand {
	return Registry
}

// Get returns a command by name
func Get(name string) (*SlashCommand, bool) {
	cmd, exists := Registry[name]
	return cmd, exists
}

// GetDefinitions returns all command definitions for registration with Discord
func GetDefinitions() []*discordgo.ApplicationCommand {
	var definitions []*discordgo.ApplicationCommand
	for _, cmd := range Registry {
		definitions = append(definitions, cmd.Definition)
	}
	return definitions
}

// Initialize registers all commands
func Initialize() {
	RegisterPing()
	RegisterHelp()
	RegisterRoll()
	RegisterSubscribe()
	RegisterUnsubscribe()
	RegisterMyGames()
	RegisterGamesList()
}
