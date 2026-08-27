package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
)

type moderationUserRepository struct {
	userrepository.UserRepository
	user *models.User
}

func (r moderationUserRepository) GetUserByToken(string) *models.User { return r.user }

var _ userrepository.UserRepository = moderationUserRepository{}

func TestRequireUserModerationScopeAttachesAuthenticatedUser(t *testing.T) {
	moderator := &models.User{
		ID:     "moderator-id",
		Scopes: []string{models.ModeratorScopeKey},
	}
	middleware := &Middleware{
		userRepository: moderationUserRepository{user: moderator},
	}
	var got *models.User
	handler := middleware.RequireUserModerationScopeAccesstoken(func(_ http.ResponseWriter, r *http.Request) {
		got = UserFromRequest(r)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/chat/messagevisibility?accessToken=token", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got == nil || got.ID != moderator.ID {
		t.Fatalf("authenticated user = %#v, want ID %q", got, moderator.ID)
	}
}
