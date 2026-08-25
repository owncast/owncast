package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	notificationsrepo "github.com/owncast/owncast/persistence/notificationsrepository"
)

type notificationRepositoryStub struct {
	addedChannel     string
	addedDestination string
}

func (r *notificationRepositoryStub) AddNotification(channel, destination string) error {
	r.addedChannel = channel
	r.addedDestination = destination
	return nil
}

func (*notificationRepositoryStub) RemoveNotificationForChannel(string, string) error {
	return nil
}

func (*notificationRepositoryStub) GetNotificationDestinationsForChannel(string) ([]string, error) {
	return nil, nil
}

func (*notificationRepositoryStub) Setup() {}

func TestRegisterForLiveNotificationsValidatesBrowserEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		wantStored  bool
	}{
		{
			name:        "public HTTPS endpoint",
			destination: `{"endpoint":"https://93.184.216.34/push","keys":{"auth":"auth","p256dh":"key"}}`,
			wantStored:  true,
		},
		{
			name:        "loopback endpoint",
			destination: `{"endpoint":"https://127.0.0.1/push"}`,
		},
		{
			name:        "link local endpoint",
			destination: `{"endpoint":"https://169.254.169.254/latest"}`,
		},
		{
			name:        "plain HTTP endpoint",
			destination: `{"endpoint":"http://93.184.216.34/push"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &notificationRepositoryStub{}
			h := &Handlers{notificationsRepository: repo}
			body := `{"channel":"` + notificationsrepo.BrowserPushNotification + `","destination":` + strconv.Quote(test.destination) + `}`
			req := httptest.NewRequest(http.MethodPost, "/api/notifications/register", strings.NewReader(body))
			rec := httptest.NewRecorder()

			h.RegisterForLiveNotifications(models.User{}, rec, req)

			stored := repo.addedDestination != ""
			if stored != test.wantStored {
				t.Fatalf("stored = %v, want %v, response = %s", stored, test.wantStored, rec.Body.String())
			}
		})
	}
}
