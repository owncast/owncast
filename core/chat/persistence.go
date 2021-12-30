package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

const (
	maxBacklogHours  = 5  // Keep backlog max hours worth of messages
	maxBacklogNumber = 50 // Return max number of messages in history request
)

func setupPersistence() {
	_datastore = data.GetDatastore()
	data.CreateMessagesTable(_datastore.DB)

	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		runPruner()
		for range chatDataPruner.C {
			runPruner()
		}
	}()
}

// SaveUserMessage will save a single chat event to the messages database.
func SaveUserMessage(event events.UserMessageEvent) {
	saveEvent(event.ID, &event.User.ID, event.Body, event.Type, event.HiddenAt, event.Timestamp, nil, nil, nil, nil)
}

func saveFederatedAction(event events.FediverseEngagementEvent) {
	saveEvent(event.ID, nil, event.Body, event.Type, nil, event.Timestamp, event.Image, &event.Link, &event.UserAccountName, nil)
}

// nolint: unparam
func saveEvent(id string, userID *string, body string, eventType string, hidden *time.Time, timestamp time.Time, image *string, link *string, title *string, subtitle *string) {
	defer func() {
		_historyCache = nil
	}()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Errorln("error saving", eventType, err)
		return
	}

	defer tx.Rollback() // nolint

	stmt, err := tx.Prepare("INSERT INTO messages(id, user_id, body, eventType, hidden_at, timestamp, image, link, title, subtitle) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Errorln("error saving", eventType, err)
		return
	}

	defer stmt.Close()

	if _, err = stmt.Exec(id, userID, body, eventType, hidden, timestamp, image, link, title, subtitle); err != nil {
		log.Errorln("error saving", eventType, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Errorln("error saving", eventType, err)
		return
	}
}

func makeUserMessageEventFromRowData(row rowData) events.UserMessageEvent {
	scopes := ""
	if row.userScopes != nil {
		scopes = *row.userScopes
	}

	previousUsernames := []string{}
	if row.previousUsernames != nil {
		previousUsernames = strings.Split(*row.previousUsernames, ",")
	}

	displayName := ""
	if row.userDisplayName != nil {
		displayName = *row.userDisplayName
	}

	displayColor := 0
	if row.userDisplayColor != nil {
		displayColor = *row.userDisplayColor
	}

	createdAt := time.Time{}
	if row.userCreatedAt != nil {
		createdAt = *row.userCreatedAt
	}

	u := user.User{
		ID:            *row.userID,
		AccessToken:   "",
		DisplayName:   displayName,
		DisplayColor:  displayColor,
		CreatedAt:     createdAt,
		DisabledAt:    row.userDisabledAt,
		NameChangedAt: row.userNameChangedAt,
		PreviousNames: previousUsernames,
		Scopes:        strings.Split(scopes, ","),
	}

	message := events.UserMessageEvent{
		Event: events.Event{
			Type:      row.eventType,
			ID:        row.id,
			Timestamp: row.timestamp,
		},
		UserEvent: events.UserEvent{
			User:     &u,
			HiddenAt: row.hiddenAt,
		},
		MessageEvent: events.MessageEvent{
			Body:    row.body,
			RawBody: row.body,
		},
	}

	return message
}

func makeSystemMessageChatEventFromRowData(row rowData) events.SystemMessageEvent {
	message := events.SystemMessageEvent{
		Event: events.Event{
			Type:      row.eventType,
			ID:        row.id,
			Timestamp: row.timestamp,
		},
		MessageEvent: events.MessageEvent{
			Body:    row.body,
			RawBody: row.body,
		},
	}
	return message
}

func makeActionMessageChatEventFromRowData(row rowData) events.ActionEvent {
	message := events.ActionEvent{
		Event: events.Event{
			Type:      row.eventType,
			ID:        row.id,
			Timestamp: row.timestamp,
		},
		MessageEvent: events.MessageEvent{
			Body:    row.body,
			RawBody: row.body,
		},
	}
	return message
}

func makeFederatedActionChatEventFromRowData(row rowData) events.FediverseEngagementEvent {
	message := events.FediverseEngagementEvent{
		Event: events.Event{
			Type:      row.eventType,
			ID:        row.id,
			Timestamp: row.timestamp,
		},
		MessageEvent: events.MessageEvent{
			Body:    row.body,
			RawBody: row.body,
		},
		Image:           row.image,
		Link:            *row.link,
		UserAccountName: *row.title,
	}
	return message
}

type rowData struct {
	id        string
	userID    *string
	body      string
	eventType models.EventType
	hiddenAt  *time.Time
	timestamp time.Time
	title     *string
	subtitle  *string
	image     *string
	link      *string

	userDisplayName   *string
	userDisplayColor  *int
	userCreatedAt     *time.Time
	userDisabledAt    *time.Time
	previousUsernames *string
	userNameChangedAt *time.Time
	userScopes        *string
}

func getChat(query string) []interface{} {
	history := make([]interface{}, 0)
	rows, err := _datastore.DB.Query(query)
	if err != nil || rows.Err() != nil {
		log.Errorln("error fetching chat history", err)
		return history
	}
	defer rows.Close()

	for rows.Next() {
		row := rowData{}

		// Convert a database row into a chat event
		if err = rows.Scan(
			&row.id,
			&row.userID,
			&row.body,
			&row.title,
			&row.subtitle,
			&row.image,
			&row.link,
			&row.eventType,
			&row.hiddenAt,
			&row.timestamp,
			&row.userDisplayName,
			&row.userDisplayColor,
			&row.userCreatedAt,
			&row.userDisabledAt,
			&row.previousUsernames,
			&row.userNameChangedAt,
			&row.userScopes,
		); err != nil {
			log.Errorln("There is a problem converting query to chat objects. Please report this:", query)
			break
		}

		var message interface{}

		switch row.eventType {
		case events.MessageSent:
			message = makeUserMessageEventFromRowData(row)
		case events.SystemMessageSent:
			message = makeSystemMessageChatEventFromRowData(row)
		case events.ChatActionSent:
			message = makeActionMessageChatEventFromRowData(row)
		case events.FediverseEngagementFollow:
			message = makeFederatedActionChatEventFromRowData(row)
		case events.FediverseEngagementLike:
			message = makeFederatedActionChatEventFromRowData(row)
		case events.FediverseEngagementRepost:
			message = makeFederatedActionChatEventFromRowData(row)
		}

		history = append(history, message)
	}

	return history
}

var _historyCache *[]interface{}

// GetChatModerationHistory will return all the chat messages suitable for moderation purposes.
func GetChatModerationHistory() []interface{} {
	if _historyCache != nil {
		return *_historyCache
	}

	// Get all messages regardless of visibility
	query := "SELECT messages.id, user_id, body, title, subtitle, image, link, eventType, hidden_at, timestamp, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, scopes FROM messages INNER JOIN users ON messages.user_id = users.id ORDER BY timestamp DESC"
	result := getChat(query)

	_historyCache = &result

	return result
}

// GetChatHistory will return all the chat messages suitable for returning as user-facing chat history.
func GetChatHistory() []interface{} {
	// Get all visible messages
	query := fmt.Sprintf("SELECT messages.id,messages.user_id, messages.body, messages.title, messages.subtitle, messages.image, messages.link, messages.eventType, messages.hidden_at, messages.timestamp, users.display_name, users.display_color, users.created_at, users.disabled_at, users.previous_names, users.namechanged_at, users.scopes FROM messages LEFT JOIN users ON messages.user_id = users.id AND hidden_at IS NULL AND disabled_at IS NULL ORDER BY timestamp DESC LIMIT %d", maxBacklogNumber)
	m := getChat(query)

	// Invert order of messages
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}

	return m
}

// SetMessageVisibilityForUserID will bulk change the visibility of messages for a user
// and then send out visibility changed events to chat clients.
func SetMessageVisibilityForUserID(userID string, visible bool) error {
	defer func() {
		_historyCache = nil
	}()

	// Get a list of IDs from this user within the 5hr window to send to the connected clients to hide
	ids := make([]string, 0)
	query := fmt.Sprintf("SELECT messages.id, user_id, body, eventType, hidden_at, timestamp, display_name, display_color, created_at, disabled_at,  previous_names, namechanged_at FROM messages INNER JOIN users ON messages.user_id = users.id WHERE user_id IS '%s'", userID)
	messages := getChat(query)

	if len(messages) == 0 {
		return nil
	}

	for _, message := range messages {
		ids = append(ids, message.(events.Event).ID)
	}

	// Tell the clients to hide/show these messages.
	return SetMessagesVisibility(ids, visible)
}

func saveMessageVisibility(messageIDs []string, visible bool) error {
	defer func() {
		_historyCache = nil
	}()

	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE messages SET hidden_at=? WHERE id IN (?" + strings.Repeat(",?", len(messageIDs)-1) + ")")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var hiddenAt *time.Time
	if !visible {
		now := time.Now()
		hiddenAt = &now
	} else {
		hiddenAt = nil
	}

	args := make([]interface{}, len(messageIDs)+1)
	args[0] = hiddenAt
	for i, id := range messageIDs {
		args[i+1] = id
	}

	if _, err = stmt.Exec(args...); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
