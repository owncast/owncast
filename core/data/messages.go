package data

import (
	"context"
	"database/sql"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// CreateMessagesTable will create the chat messages table if needed.
func CreateMessagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL,
		"user_id" INTEGER,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATETIME,
		"timestamp" DATETIME,
    "title" TEXT,
    "subtitle" TEXT,
    "image" TEXT,
    "link" TEXT,
		PRIMARY KEY (id)
	);CREATE INDEX index ON messages (id, user_id, hidden_at, timestamp);
	CREATE INDEX id ON messages (id);
	CREATE INDEX user_id ON messages (user_id);
	CREATE INDEX hidden_at ON messages (hidden_at);
	CREATE INDEX timestamp ON messages (timestamp);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("error creating chat messages table", err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal("error creating chat messages table", err)
	}
}

// GetMessagesCount will return the number of messages in the database.
func GetMessagesCount() int64 {
	query := `SELECT COUNT(*) FROM messages`
	rows, err := _db.Query(query)
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

// CreateBanIPTable will create the IP ban table if needed.
func CreateBanIPTable(db *sql.DB) {
	createTableSQL := `  CREATE TABLE IF NOT EXISTS ip_bans (
    "ip_address" TEXT NOT NULL PRIMARY KEY,
    "notes" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("error creating ip ban table", err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal("error creating ip ban table", err)
	}
}

// BanIPAddress will persist a new IP address ban to the datastore.
func BanIPAddress(address, note string) error {
	return _datastore.GetQueries().BanIPAddress(context.Background(), db.BanIPAddressParams{
		IpAddress: address,
		Notes:     sql.NullString{String: note, Valid: true},
	})
}

// IsIPAddressBanned will return if an IP address has been previously blocked.
func IsIPAddressBanned(address string) (bool, error) {
	blocked, error := _datastore.GetQueries().IsIPAddressBlocked(context.Background(), address)
	return blocked > 0, error
}

// GetIPAddressBans will return all the banned IP addresses.
func GetIPAddressBans() ([]models.IPAddress, error) {
	result, err := _datastore.GetQueries().GetIPAddressBans(context.Background())
	if err != nil {
		return nil, err
	}

	response := []models.IPAddress{}
	for _, ip := range result {
		response = append(response, models.IPAddress{
			IPAddress: ip.IpAddress,
			Notes:     ip.Notes.String,
			CreatedAt: ip.CreatedAt.Time,
		})
	}
	return response, err
}

// RemoveIPAddressBan will remove a previously banned IP address.
func RemoveIPAddressBan(address string) error {
	return _datastore.GetQueries().RemoveIPAddressBan(context.Background(), address)
}
