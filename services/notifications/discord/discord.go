package discord

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Discord is an instance of the Discord service.
type Discord struct {
	name       string
	avatar     string
	webhookURL string
}

// New will create a new instance of the Discord service.
func New(name, avatar, webhook string) (*Discord, error) {
	return &Discord{
		name:       name,
		avatar:     avatar,
		webhookURL: webhook,
	}, nil
}

// Send will send a message to a Discord channel via a webhook.
func (t *Discord) Send(content string) error {
	type message struct {
		Username string `json:"username"`
		Content  string `json:"content"`
		Avatar   string `json:"avatar_url"`
	}

	msg := message{
		Username: t.name,
		Content:  content,
		Avatar:   t.avatar,
	}

	jsonText, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "error marshalling discord message to json")
	}

	req, err := http.NewRequest("POST", t.webhookURL, bytes.NewReader(jsonText))
	if err != nil {
		return errors.Wrap(err, "error creating discord webhook request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing discord webhook")
	}

	return resp.Body.Close()
}
