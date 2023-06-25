package userrepository

import (
	"context"
	"strings"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/sqlstorage"
)

// AddAuth will add an external authentication token and type for a user.
func (r *SqlUserRepository) AddAuth(userID, authToken string, authType models.AuthType) error {
	return r.datastore.GetQueries().AddAuthForUser(context.Background(), sqlstorage.AddAuthForUserParams{
		UserID: userID,
		Token:  authToken,
		Type:   string(authType),
	})
}

// GetUserByAuth will return an existing user given auth details if a user
// has previously authenticated with that method.
func (r *SqlUserRepository) GetUserByAuth(authToken string, authType models.AuthType) *models.User {
	u, err := r.datastore.GetQueries().GetUserByAuth(context.Background(), sqlstorage.GetUserByAuthParams{
		Token: authToken,
		Type:  string(authType),
	})
	if err != nil {
		return nil
	}

	var scopes []string
	if u.Scopes.Valid {
		scopes = strings.Split(u.Scopes.String, ",")
	}

	return &models.User{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		DisplayColor:    int(u.DisplayColor),
		CreatedAt:       u.CreatedAt.Time,
		DisabledAt:      &u.DisabledAt.Time,
		PreviousNames:   strings.Split(u.PreviousNames.String, ","),
		NameChangedAt:   &u.NamechangedAt.Time,
		AuthenticatedAt: &u.AuthenticatedAt.Time,
		Scopes:          scopes,
	}
}
