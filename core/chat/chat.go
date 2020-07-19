package chat

import (
	"errors"
	"time"

	"github.com/gabek/owncast/models"
)

//Setup sets up the chat server
func Setup(listener models.ChatListener) {
	setupPersistence()

	messages := []models.ChatMessage{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan models.ChatMessage)
	pingCh := make(chan models.PingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)
	messages = append(messages, getChatHistory()...)

	_server = &server{
		messages,
		clients,
		"/entry", //hardcoded due to the UI requiring this and it is not configurable
		listener,
		addCh,
		delCh,
		sendAllCh,
		pingCh,
		doneCh,
		errCh,
	}
}

//Start starts the chat server
func Start() error {
	if _server == nil {
		return errors.New("chat server is nil")
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				_server.ping()
			}
		}
	}()

	_server.Listen()

	return errors.New("chat server failed to start")
}

//SendMessage sends a message to all
func SendMessage(message models.ChatMessage) {
	if _server == nil {
		return
	}

	_server.SendToAll(message)
}

//GetMessages gets all of the messages
func GetMessages() []models.ChatMessage {
	if _server == nil {
		return []models.ChatMessage{}
	}

	return getRecentMessages()
}

func getRecentMessages() []models.ChatMessage {
	if len(_server.Messages) < 100 {
		return _server.Messages
	}

	maxAgeInHours := float64(2)
	messages := make([]models.ChatMessage, 0)
	for _, message := range _server.Messages {
		if time.Since(message.Timestamp).Hours() < maxAgeInHours {
			messages = append(messages, message)
		}
	}

	return messages
}
