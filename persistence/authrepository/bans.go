package authrepository

import (
	"context"
	"database/sql"
	"log"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/webserver/handlers/generated"
)

// CreateBanIPTable will create the IP ban table if needed.
func (r *SqlAuthRepository) CreateBanIPTable(db *sql.DB) {
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
func (r *SqlAuthRepository) BanIPAddress(address, note string) error {
	return r.datastore.GetQueries().BanIPAddress(context.Background(), db.BanIPAddressParams{
		IpAddress: address,
		Notes:     sql.NullString{String: note, Valid: true},
	})
}

// IsIPAddressBanned will return if an IP address has been previously blocked.
func (r *SqlAuthRepository) IsIPAddressBanned(address string) (bool, error) {
	blocked, error := r.datastore.GetQueries().IsIPAddressBlocked(context.Background(), address)
	return blocked > 0, error
}

// GetIPAddressBans will return all the banned IP addresses.
func (r *SqlAuthRepository) GetIPAddressBans() ([]generated.IPAddress, error) {
	result, err := r.datastore.GetQueries().GetIPAddressBans(context.Background())
	if err != nil {
		return nil, err
	}

	response := []generated.IPAddress{}
	for _, ip := range result {
		response = append(response, generated.IPAddress{
			IpAddress: &ip.IpAddress,
			Notes:     &ip.Notes.String,
			CreatedAt: &ip.CreatedAt.Time,
		})
	}
	return response, err
}

// RemoveIPAddressBan will remove a previously banned IP address.
func (r *SqlAuthRepository) RemoveIPAddressBan(address string) error {
	return r.datastore.GetQueries().RemoveIPAddressBan(context.Background(), address)
}
