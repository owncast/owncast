package persistence

import (
	"context"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

// Setup will initialize the ActivityPub persistence layer with the provided datastore.
func Setup(datastore *data.Datastore) {
	_datastore = datastore
	createFederationFollowersTable()
	createFederationOutboxTable()
}

// AddFollow will save a follow to the datastore.
func AddFollow(follow apmodels.ActivityPubActor) error {
	log.Println("Saving", follow.ActorIri, "as a follower.")
	return createFollow(follow.ActorIri.String(), follow.Inbox.String(), follow.Name, follow.Username, follow.Image.String())
}

// RemoveFollow will remove a follow from the datastore.
func RemoveFollow(unfollow apmodels.ActivityPubActor) error {
	log.Println("Removing", unfollow.ActorIri, "as a follower.")
	return removeFollow(unfollow.ActorIri)
}

func createFollow(actor string, inbox string, name string, username string, image string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO ap_followers(iri, inbox, name, username, image) values(?, ?, ?, ?, ?)")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(actor, inbox, name, username, image)
	if err != nil {
		log.Errorln("error creating new federation follow", err)
	}

	return tx.Commit()
}

// UpdateFollower will update the details of a stored follower given an IRI.
func UpdateFollower(actorIRI string, inbox string, name string, username string, image string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("UPDATE ap_followers SET inbox = ?, name = ?, username = ?, image = ? WHERE iri IS ?")
	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(actorIRI)
	if err != nil {
		log.Errorln("error removing federation follow", err)
	}

	return tx.Commit()
}

func removeFollow(actor *url.URL) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("DELETE FROM ap_followers WHERE iri IS ?")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(actor.String())
	if err != nil {
		log.Errorln("error removing federation follow", err)
	}

	return tx.Commit()
}

func createFederationOutboxTable() {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_outbox (
		"id" TEXT NOT NULL,
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id));
		CREATE INDEX id ON ap_outbox (id);
		CREATE INDEX iri ON ap_outbox (iri);
		CREATE INDEX type ON ap_outbox (type);
	);`

	stmt, err := _datastore.DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln("error executing sql creating outbox table", createTableSQL, err)
	}
}

func createFederationFollowersTable() {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (iri));
		CREATE INDEX iri ON ap_followers (iri);
	);`

	stmt, err := _datastore.DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln("error executing sql creating followers table", createTableSQL, err)
	}
}

// GetOutbox will create an instance of the outbox populated by stored items.
func GetOutbox() (vocab.ActivityStreamsOrderedCollectionPage, error) {
	collection := streams.NewActivityStreamsOrderedCollectionPage()
	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()

	query := `SELECT value FROM ap_outbox`

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return collection, err
	}
	defer rows.Close()

	for rows.Next() {
		var value []byte
		if err := rows.Scan(&value); err != nil {
			log.Error("There is a problem reading the database.", err)
			return collection, err
		}

		createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
			// items = append(items, activity)
			orderedItems.AppendActivityStreamsCreate(activity)
			return nil
		}

		if err := resolvers.Resolve(context.TODO(), value, createCallback); err != nil {
			return collection, err
		}
	}

	if err := rows.Err(); err != nil {
		return collection, err
	}

	collection.SetActivityStreamsOrderedItems(orderedItems)

	query = `SElECT count(*) FROM ap_outbox`
	rows, err = _datastore.DB.Query(query)
	if err != nil {
		return collection, err
	}
	defer rows.Close()
	rows.Next()
	var totalCount int
	if err := rows.Scan(&totalCount); err != nil {
		return collection, err
	}

	totalItems := streams.NewActivityStreamsTotalItemsProperty()
	totalItems.Set(totalCount)
	collection.SetActivityStreamsTotalItems(totalItems)

	return collection, nil
}

// AddToOutbox will store a single payload to the persistence layer.
func AddToOutbox(id string, iri string, itemData []byte, typeString string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO ap_outbox(id, iri, value, type) values(?, ?, ?, ?)")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, iri, itemData, typeString)
	if err != nil {
		log.Errorln("error creating new item in federation outbox", err)
	}

	return tx.Commit()
}

// GetObjectByID will return a string representation of a single object by the ID.
func GetObjectByID(id string) (string, error) {
	query := `SELECT value FROM ap_outbox WHERE id IS ?`
	row := _datastore.DB.QueryRow(query, id)

	var value string
	err := row.Scan(&value)

	return value, err
}

// GetObjectByIRI will return a string representation of a single object by the IRI.
func GetObjectByIRI(IRI string) (string, error) {
	query := `SELECT value FROM ap_outbox WHERE iri IS ?`
	// log.Println(query, IRI)
	row := _datastore.DB.QueryRow(query, IRI)

	var value string
	err := row.Scan(&value)

	return value, err
}

// GetLocalPostCount will return the number of posts existing locally.
func GetLocalPostCount() int {
	var totalCount int

	query := `SElECT count(*) FROM ap_outbox`
	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return totalCount
	}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&totalCount); err != nil {
		return 0
	}

	return totalCount
}

// GetFollowerCount will return the number of followers we're keeping track of.
func GetFollowerCount() int {
	var totalCount int

	query := `SElECT count(*) FROM ap_followers`
	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return totalCount
	}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&totalCount); err != nil {
		return 0
	}

	return totalCount
}

// GetFederationFollowers will return a slice of the followers we keep track of locally.
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
