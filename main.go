package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/bot"
	"discord-bot/commands"
	"discord-bot/config"
	"discord-bot/data"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize subscription manager
	commands.SubManager = data.NewSubscriptionManager("subscriptions.json")

	// Create bot
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	// Start the bot
	err = b.Start()
	if err != nil {
		log.Fatal("Error starting bot: ", err)
	}
	defer b.Stop()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait for a CTRL+C signal to gracefully shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
