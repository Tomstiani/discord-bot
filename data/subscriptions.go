package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// GameSubscription represents a user's subscription to a game
type GameSubscription struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Game      string `json:"game"`
	NTFYTopic string `json:"ntfy_topic"` // Their personal NTFY topic
}

// SubscriptionManager manages game subscriptions
type SubscriptionManager struct {
	subscriptions []GameSubscription
	filePath      string
	mutex         sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(filePath string) *SubscriptionManager {
	sm := &SubscriptionManager{
		subscriptions: make([]GameSubscription, 0),
		filePath:      filePath,
	}
	sm.loadFromFile()
	return sm
}

// Subscribe adds a user's subscription to a game
func (sm *SubscriptionManager) Subscribe(userID, username, game, ntfyTopic string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if already subscribed
	for _, sub := range sm.subscriptions {
		if sub.UserID == userID && sub.Game == game {
			return fmt.Errorf("already subscribed to %s", game)
		}
	}

	// Add new subscription
	newSub := GameSubscription{
		UserID:    userID,
		Username:  username,
		Game:      game,
		NTFYTopic: ntfyTopic,
	}
	sm.subscriptions = append(sm.subscriptions, newSub)

	return sm.saveToFile()
}

// Unsubscribe removes a user's subscription to a game
func (sm *SubscriptionManager) Unsubscribe(userID, game string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i, sub := range sm.subscriptions {
		if sub.UserID == userID && sub.Game == game {
			// Remove subscription
			sm.subscriptions = append(sm.subscriptions[:i], sm.subscriptions[i+1:]...)
			return sm.saveToFile()
		}
	}

	return fmt.Errorf("not subscribed to %s", game)
}

// GetSubscriptions returns all subscriptions for a user
func (sm *SubscriptionManager) GetSubscriptions(userID string) []GameSubscription {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var userSubs []GameSubscription
	for _, sub := range sm.subscriptions {
		if sub.UserID == userID {
			userSubs = append(userSubs, sub)
		}
	}
	return userSubs
}

// GetSubscribersForGame returns all subscribers for a specific game
func (sm *SubscriptionManager) GetSubscribersForGame(game string) []GameSubscription {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var subscribers []GameSubscription
	for _, sub := range sm.subscriptions {
		if sub.Game == game {
			subscribers = append(subscribers, sub)
		}
	}
	return subscribers
}

// GetAllGames returns a list of all games people are subscribed to
func (sm *SubscriptionManager) GetAllGames() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	gameSet := make(map[string]bool)
	for _, sub := range sm.subscriptions {
		gameSet[sub.Game] = true
	}

	games := make([]string, 0, len(gameSet))
	for game := range gameSet {
		games = append(games, game)
	}
	return games
}

// saveToFile saves subscriptions to JSON file
func (sm *SubscriptionManager) saveToFile() error {
	data, err := json.MarshalIndent(sm.subscriptions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.filePath, data, 0644)
}

// loadFromFile loads subscriptions from JSON file
func (sm *SubscriptionManager) loadFromFile() error {
	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &sm.subscriptions)
}
