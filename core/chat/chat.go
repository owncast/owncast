package chat

import (
	"errors"
	"time"

	"github.com/gabek/owncast/models"
)

//Setup sets up the chat server
func Setup(listener models.ChatListener) {
	messages := []models.ChatMessage{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan models.ChatMessage)
	pingCh := make(chan models.PingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	// Demo messages only.  Remove me eventually!!!
	messages = append(messages, models.ChatMessage{"", "Tom Nook", "I'll be there with Bells on! Ho ho!", "https://gamepedia.cursecdn.com/animalcrossingpocketcamp_gamepedia_en/thumb/4/4f/Timmy_Icon.png/120px-Timmy_Icon.png?version=87b38d7d6130411d113486c2db151385", "demo-message-1", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"", "Redd", "Fool me once, shame on you. Fool me twice, stop foolin' me.", "https://vignette.wikia.nocookie.net/animalcrossing/images/3/3d/Redd2.gif/revision/latest?cb=20100710004252", "demo-message-2", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"", "Kevin", "You just caught me before I was about to go work out weeweewee!", "https://vignette.wikia.nocookie.net/animalcrossing/images/2/20/NH-Kevin_poster.png/revision/latest/scale-to-width-down/100?cb=20200410185817", "demo-message-3", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"", "Isabelle", " Isabelle is the mayor's highly capable secretary. She can be forgetful sometimes, but you can always count on her for information about the town. She wears her hair up in a bun that makes her look like a shih tzu. Mostly because she is one! She also has a twin brother named Digby.", "https://dodo.ac/np/images/thumb/7/7b/IsabelleTrophyWiiU.png/200px-IsabelleTrophyWiiU.png", "demo-message-4", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"", "Judy", "myohmy, I'm dancing my dreams away.", "https://vignette.wikia.nocookie.net/animalcrossing/images/5/50/NH-Judy_poster.png/revision/latest/scale-to-width-down/100?cb=20200522063219", "demo-message-5", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"", "Blathers", "Blathers is an owl with brown feathers. His face is white and he has a yellow beak. His arms are wing shaped and he has yellow talons. His eyes are very big with small black irises. He also has big pink cheek circles on his cheeks. His belly appears to be checkered in diamonds with light brown and white squares, similar to an argyle vest, which is traditionally associated with academia. His green bowtie further alludes to his academic nature.", "https://vignette.wikia.nocookie.net/animalcrossing/images/b/b3/NH-character-Blathers.png/revision/latest?cb=20200229053519", "demo-message-6", "ChatMessage"})

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

	return _server.Messages
}
