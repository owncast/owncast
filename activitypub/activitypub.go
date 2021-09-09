package activitypub

import (
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

	var query = "SELECT iri, inbox, name, username, image FROM ap_followers"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return followers, err
	}
	defer rows.Close()

	for rows.Next() {
		var iriString string
		var inboxString string
		var nameString string
		var usernameString string
		var imageString string

		if err := rows.Scan(&iriString, &inboxString, &nameString, &usernameString, &imageString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return followers, err
		}

		singleFollower := models.Follower{
			Name:     nameString,
			Username: usernameString,
			Image:    imageString,
			Link:     iriString,
			Inbox:    inboxString,
		}

		followers = append(followers, singleFollower)
	}

	if err := rows.Err(); err != nil {
		return followers, err
	}

	return followers, nil
}
