package chatmessagerepository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"

	log "github.com/sirupsen/logrus"
)

const maxBacklogNumber = 50 // Return max number of messages in history request

type ChatMessageRepository interface {
	SaveUserMessage(event events.UserMessageEvent)
	SaveFederatedAction(event events.FediverseEngagementEvent)
	SaveEvent(id string, userID *string, body string, eventType string, hidden *time.Time, timestamp time.Time, image *string, link *string, title *string, subtitle *string)
	GetChatModerationHistory() []interface{}
	GetChatHistory() []interface{}
	GetMessagesFromUser(userID string) ([]events.UserMessageEvent, error)
	GetMessageIdsForUserID(userID string) ([]string, error)
	SetMessageVisibilityForMessageIDs(messageIDs []string, visible bool) error
	GetMessagesCount() int64
}

type SqlChatMessageRepository struct {
	datastore *data.Datastore
}

// NOTE: This is temporary during the transition period.
var temporaryGlobalInstance ChatMessageRepository

// Get will return the user repository.
func Get() ChatMessageRepository {
	if temporaryGlobalInstance == nil {
		i := New(data.GetDatastore())
		temporaryGlobalInstance = i
	}
	return temporaryGlobalInstance
}

// New will create a new instance of the UserRepository.
func New(datastore *data.Datastore) ChatMessageRepository {
	r := SqlChatMessageRepository{
		datastore: datastore,
	}

	return &r
}

// SaveUserMessage will save a single chat event to the messages database.
func (r *SqlChatMessageRepository) SaveUserMessage(event events.UserMessageEvent) {
	r.SaveEvent(event.ID, &event.User.ID, event.Body, event.Type, event.HiddenAt, event.Timestamp, nil, nil, nil, nil)
}

func (r *SqlChatMessageRepository) SaveFederatedAction(event events.FediverseEngagementEvent) {
	r.SaveEvent(event.ID, nil, event.Body, event.Type, nil, event.Timestamp, event.Image, &event.Link, &event.UserAccountName, nil)
}

// nolint: unparam
func (r *SqlChatMessageRepository) SaveEvent(id string, userID *string, body string, eventType string, hidden *time.Time, timestamp time.Time, image *string, link *string, title *string, subtitle *string) {
	defer func() {
		_historyCache = nil
	}()

	tx, err := r.datastore.DB.Begin()
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

	isBot := (row.userType != nil && *row.userType == "API")
	scopeSlice := strings.Split(scopes, ",")

	u := models.User{
		ID:              *row.userID,
		DisplayName:     displayName,
		DisplayColor:    displayColor,
		CreatedAt:       createdAt,
		DisabledAt:      row.userDisabledAt,
		NameChangedAt:   row.userNameChangedAt,
		PreviousNames:   previousUsernames,
		AuthenticatedAt: row.userAuthenticatedAt,
		Authenticated:   row.userAuthenticatedAt != nil,
		Scopes:          scopeSlice,
		IsBot:           isBot,
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
	timestamp         time.Time
	image             *string
	previousUsernames *string

	userDisplayName  *string
	userDisplayColor *int
	userID           *string
	title            *string
	subtitle         *string
	link             *string

	userType            *string
	userScopes          *string
	hiddenAt            *time.Time
	userCreatedAt       *time.Time
	userDisabledAt      *time.Time
	userAuthenticatedAt *time.Time
	userNameChangedAt   *time.Time
	body                string
	eventType           models.EventType
	id                  string
}

func getChat(rows *sql.Rows) ([]interface{}, error) {
	history := make([]interface{}, 0)

	for rows.Next() {
		row := rowData{}

		// Convert a database row into a chat event
		if err := rows.Scan(
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
			&row.userAuthenticatedAt,
			&row.userScopes,
			&row.userType,
		); err != nil {
			return nil, err
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

	return history, nil
}

var _historyCache *[]interface{}

// GetChatModerationHistory will return all the chat messages suitable for moderation purposes.
func (r *SqlChatMessageRepository) GetChatModerationHistory() []interface{} {
	if _historyCache != nil {
		return *_historyCache
	}

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		log.Errorln("error fetching chat moderation history", err)
		return nil
	}

	defer tx.Rollback() // nolint

	// Get all messages regardless of visibility
	query := "SELECT messages.id, user_id, body, title, subtitle, image, link, eventType, hidden_at, timestamp, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes, type FROM messages INNER JOIN users ON messages.user_id = users.id ORDER BY timestamp DESC"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Errorln("error fetching chat moderation history", err)
		return nil
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Errorln("error fetching chat moderation history", err)
		return nil
	}

	defer stmt.Close()
	defer rows.Close()

	result, err := getChat(rows)
	if err != nil {
		log.Errorln(err)
		log.Errorln("There is a problem enumerating chat message rows. Please report this:", query)
		return nil
	}

	_historyCache = &result

	if err = tx.Commit(); err != nil {
		log.Errorln("error fetching chat moderation history", err)
		return nil
	}

	return result
}

// GetChatHistory will return all the chat messages suitable for returning as user-facing chat history.
func (r *SqlChatMessageRepository) GetChatHistory() []interface{} {
	tx, err := r.datastore.DB.Begin()
	if err != nil {
		log.Errorln("error fetching chat history", err)
		return nil
	}

	defer tx.Rollback() // nolint

	// Get all visible messages
	query := "SELECT messages.id, messages.user_id, messages.body, messages.title, messages.subtitle, messages.image, messages.link, messages.eventType, messages.hidden_at, messages.timestamp, users.display_name, users.display_color, users.created_at, users.disabled_at, users.previous_names, users.namechanged_at, users.authenticated_at, users.scopes, users.type FROM users JOIN messages ON users.id = messages.user_id WHERE hidden_at IS NULL AND disabled_at IS NULL ORDER BY timestamp DESC LIMIT ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Errorln("error fetching chat history", err)
		return nil
	}

	rows, err := stmt.Query(maxBacklogNumber)
	if err != nil {
		log.Errorln("error fetching chat history", err)
		return nil
	}

	defer stmt.Close()
	defer rows.Close()

	m, err := getChat(rows)
	if err != nil {
		log.Errorln(err)
		log.Errorln("There is a problem enumerating chat message rows. Please report this:", query)
		return nil
	}

	if err = tx.Commit(); err != nil {
		log.Errorln("error fetching chat history", err)
		return nil
	}

	// Invert order of messages
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}

	return m
}

// GetMessagesFromUser returns chat messages that were sent by a specific user.
func (r *SqlChatMessageRepository) GetMessagesFromUser(userID string) ([]events.UserMessageEvent, error) {
	query, err := r.datastore.GetQueries().GetMessagesFromUser(context.Background(), sql.NullString{String: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	results := make([]events.UserMessageEvent, len(query))
	for i, row := range query {
		results[i] = events.UserMessageEvent{
			Event: events.Event{
				Timestamp: row.Timestamp.Time,
				ID:        row.ID,
			},
			MessageEvent: events.MessageEvent{
				Body: row.Body.String,
			},
		}
	}

	return results, nil
}

// GetMessageIdsForUserID will return the chat message IDs for a specific user.
func (r *SqlChatMessageRepository) GetMessageIdsForUserID(userID string) ([]string, error) {
	defer func() {
		_historyCache = nil
	}()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return nil, errors.New("error while setting message visibility")
	}

	defer tx.Rollback() // nolint
	query := "SELECT messages.id, user_id, body, title, subtitle, image, link, eventType, hidden_at, timestamp, display_name, display_color, created_at, disabled_at,  previous_names, namechanged_at, authenticated_at, scopes, type FROM messages INNER JOIN users ON messages.user_id = users.id WHERE user_id IS ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, errors.New("error while setting message visibility")
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Errorln("error while setting message visibility", err)
		return nil, errors.New("error while setting message visibility")
	}

	defer stmt.Close()
	defer rows.Close()

	// Get a list of IDs to send to the connected clients to hide
	ids := make([]string, 0)

	messages, err := getChat(rows)
	if err != nil {
		log.Errorln(err)
		return nil, errors.New("There is a problem enumerating chat message rows. Please report this: " + query)
	}

	if len(messages) == 0 {
		return []string{}, nil
	}

	for _, message := range messages {
		ids = append(ids, message.(events.UserMessageEvent).ID)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.New("error while setting message visibility")
	}

	// Tell the clients to hide/show these messages.
	return ids, nil
}

func (r *SqlChatMessageRepository) SetMessageVisibilityForMessageIDs(messageIDs []string, visible bool) error {
	defer func() {
		_historyCache = nil
	}()

	if len(messageIDs) == 0 {
		return nil
	}

	r.datastore.DbLock.Lock()
	defer r.datastore.DbLock.Unlock()

	tx, err := r.datastore.DB.Begin()
	if err != nil {
		return err
	}

	// nolint:gosec
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

// GetMessagesCount will return the number of messages in the database.
func (r *SqlChatMessageRepository) GetMessagesCount() int64 {
	query := `SELECT COUNT(*) FROM messages`
	rows, err := r.datastore.DB.Query(query)
	if err != nil || rows.Err() != nil {
		return 0
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0
		}
	}
	return count
}
