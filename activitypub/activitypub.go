package activitypub

import (
	"net/url"

	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

func Start(datastore *data.Datastore) {
	_datastore = datastore

	persistence.Setup(datastore)
	StartRouter()

	// Test
	if data.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = data.SetPrivateKey(string(privateKey))
		_ = data.SetPublicKey(string(publicKey))
		if err != nil {
			log.Errorln("Unable to get private key", err)
		}
	}
}

func SendLive() {
	outbox.SendLive()
}

func SendPublicFederatedMessage(message string) {
	outbox.SendPublicMessage(message)
}

func GetFollowerCount() int {
	return persistence.GetFollowerCount()
}

func GetFederationFollowers() ([]models.Follower, error) {
	followers := make([]models.Follower, 0)

	var query = "SELECT iri, inbox FROM ap_followers"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return followers, err
	}
	defer rows.Close()

	for rows.Next() {
		var iriString string
		var inboxString string

		if err := rows.Scan(&iriString, &inboxString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return followers, err
		}

		iri, _ := url.Parse(iriString)
		inbox, _ := url.Parse(inboxString)

		singleFollower := models.Follower{
			Name:  "FILL IN NAME HERE",
			Image: "",
			Link:  iri.String(),
			Inbox: *inbox,
		}

		followers = append(followers, singleFollower)
	}

	if err := rows.Err(); err != nil {
		return followers, err
	}

	return followers, nil
}
