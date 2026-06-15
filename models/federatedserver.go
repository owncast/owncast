package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/owncast/owncast/db"
)

// FederatedServer represents a federated Owncast server that we follow.
type FederatedServer struct {
	ID                int64      `json:"id"`
	IRI               string     `json:"iri"`
	Name              *string    `json:"name,omitempty"`
	LogoURL           *string    `json:"logoUrl,omitempty"`
	IsOnline          bool       `json:"isOnline"`
	StreamTitle       *string    `json:"streamTitle,omitempty"`
	StreamDescription *string    `json:"streamDescription,omitempty"`
	Tags              []string   `json:"tags,omitempty"`
	ThumbnailURL      *string    `json:"thumbnailUrl,omitempty"`
	LastSeenOnline    *time.Time `json:"lastSeenOnline,omitempty"`
	LastStatusUpdate  *time.Time `json:"lastStatusUpdate,omitempty"`
	AddedAt           time.Time  `json:"addedAt"`
	FollowedAt        *time.Time `json:"followedAt,omitempty"`
	Pending           bool       `json:"pending"`
	Username          *string    `json:"username,omitempty"`
	DisplayName       *string    `json:"displayName,omitempty"`
	Summary           *string    `json:"summary,omitempty"`
	AcceptedAt        *time.Time `json:"acceptedAt,omitempty"`
	RejectedAt        *time.Time `json:"rejectedAt,omitempty"`
	FollowStatus      string     `json:"followStatus"`
}

// FederatedStreamUpdate represents stream metadata for ActivityPub handlers.
// Reuses existing patterns from models.Status.
type FederatedStreamUpdate struct {
	Title        *string  `json:"title,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	ThumbnailURL *string  `json:"thumbnailUrl,omitempty"`
}

// FromDatabaseModel converts a database FederatedServer model to the API model.
func (f *FederatedServer) FromDatabaseModel(dbServer db.FederatedServer) {
	f.ID = dbServer.ID
	f.IRI = dbServer.Iri
	f.Name = nullStringToPointer(dbServer.Name)
	f.LogoURL = nullStringToPointer(dbServer.LogoUrl)
	f.IsOnline = dbServer.IsOnline.Bool
	f.StreamTitle = nullStringToPointer(dbServer.StreamTitle)
	f.StreamDescription = nullStringToPointer(dbServer.StreamDescription)
	f.ThumbnailURL = nullStringToPointer(dbServer.ThumbnailUrl)
	f.LastSeenOnline = nullTimeToPointer(dbServer.LastSeenOnline)
	f.LastStatusUpdate = nullTimeToPointer(dbServer.LastStatusUpdate)
	f.AddedAt = dbServer.AddedAt.Time
	f.FollowedAt = nullTimeToPointer(dbServer.FollowedAt)
	f.Pending = dbServer.Pending.Bool
	f.Username = nullStringToPointer(dbServer.Username)
	f.DisplayName = nullStringToPointer(dbServer.DisplayName)
	f.Summary = nullStringToPointer(dbServer.Summary)
	f.AcceptedAt = nullTimeToPointer(dbServer.AcceptedAt)
	f.RejectedAt = nullTimeToPointer(dbServer.RejectedAt)

	// Default follow status to "pending" if not set
	if dbServer.FollowStatus.Valid {
		f.FollowStatus = dbServer.FollowStatus.String
	} else {
		f.FollowStatus = "pending"
	}

	// Parse tags from JSON string
	if dbServer.StreamTags.Valid && dbServer.StreamTags.String != "" {
		var tags []string
		if err := json.Unmarshal([]byte(dbServer.StreamTags.String), &tags); err == nil {
			f.Tags = tags
		}
	}
}

// Helper functions for null value conversion.
func nullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPointer(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// PointerToNullString converts a string pointer to sql.NullString.
func PointerToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

// PointerToNullTime converts a time pointer to sql.NullTime.
func PointerToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

// StringSliceToNullString converts a string slice to a JSON-encoded sql.NullString.
func StringSliceToNullString(tags []string) sql.NullString {
	if len(tags) > 0 {
		if jsonData, err := json.Marshal(tags); err == nil {
			return sql.NullString{String: string(jsonData), Valid: true}
		}
	}
	return sql.NullString{}
}

// BoolToNullBool converts a bool to sql.NullBool.
func BoolToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// TimeToNullTime converts a time.Time to sql.NullTime.
func TimeToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}
