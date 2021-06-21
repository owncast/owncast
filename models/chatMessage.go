package models

// // ChatEvent represents a single chat message.
// type ChatEvent struct {
// 	ClientID string `json:"-"`

// 	Author      string    `json:"author,omitempty"`
// 	Body        string    `json:"body,omitempty"`
// 	RawBody     string    `json:"-"`
// 	ID          string    `json:"id"`
// 	MessageType EventType `json:"type"`
// 	Visible     bool      `json:"visible"`
// 	Timestamp   time.Time `json:"timestamp,omitempty"`
// 	Ephemeral   bool      `json:"ephemeral,omitempty"`
// }

// // Valid checks to ensure the message is valid.
// func (m ChatEvent) Valid() bool {
// 	return m.Author != "" && m.Body != "" && m.ID != ""
// }

// // SetDefaults will set default values on a chat event object.
// func (m *ChatEvent) SetDefaults() {
// 	id, _ := shortid.Generate()
// 	m.ID = id
// 	m.Timestamp = time.Now()
// 	m.Visible = true
// }
