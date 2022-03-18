package browser

import (
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
)

// Browser is an instance of the Browser service.
type Browser struct {
	datastore  *data.Datastore
	privateKey string
	publicKey  string
}

// New will create a new instance of the Browser service.
func New(datastore *data.Datastore, publicKey, privateKey string) (*Browser, error) {
	return &Browser{
		datastore:  datastore,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// GenerateBrowserPushKeys will create the VAPID keys required for web push notifications.
func GenerateBrowserPushKeys() (string, string, error) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return "", "", errors.Wrap(err, "error generating web push keys")
	}

	return privateKey, publicKey, nil
}

// Send will send a browser push notification to the given subscription.
func (b *Browser) Send(
	subscription string,
	title string,
	body string,
) (bool, error) {
	type message struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Icon  string `json:"icon"`
	}

	m := message{
		Title: title,
		Body:  body,
		Icon:  "/logo/external",
	}

	d, err := json.Marshal(m)
	if err != nil {
		return false, errors.Wrap(err, "error marshalling web push message")
	}

	// Decode subscription
	s := &webpush.Subscription{}
	if err := json.Unmarshal([]byte(subscription), s); err != nil {
		return false, errors.Wrap(err, "error decoding destination subscription")
	}

	// Send Notification
	resp, err := webpush.SendNotification(d, s, &webpush.Options{
		VAPIDPublicKey:  b.publicKey,
		VAPIDPrivateKey: b.privateKey,
		Topic:           "owncast-go-live",
		TTL:             10,
		// Not really the subscriber, but a contact point for the sender.
		Subscriber: "owncast@owncast.online",
	})
	if resp.StatusCode == 410 {
		return true, nil
	} else if err != nil {
		return false, errors.Wrap(err, "error sending browser push notification")
	}
	defer resp.Body.Close()

	return false, err
}
