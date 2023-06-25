package chatrepository

import (
	"context"
	"database/sql"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/sqlstorage"
)

// GetMessagesCount will return the number of messages in the database.
func (c *ChatRepository) GetMessagesCount() int64 {
	query := `SELECT COUNT(*) FROM messages`
	rows, err := c.datastore.DB.Query(query)
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

// BanIPAddress will persist a new IP address ban to the datastore.
func (c *ChatRepository) BanIPAddress(address, note string) error {
	return c.datastore.GetQueries().BanIPAddress(context.Background(), sqlstorage.BanIPAddressParams{
		IpAddress: address,
		Notes:     sql.NullString{String: note, Valid: true},
	})
}

// IsIPAddressBanned will return if an IP address has been previously blocked.
func (c *ChatRepository) IsIPAddressBanned(address string) (bool, error) {
	blocked, error := c.datastore.GetQueries().IsIPAddressBlocked(context.Background(), address)
	return blocked > 0, error
}

// GetIPAddressBans will return all the banned IP addresses.
func (c *ChatRepository) GetIPAddressBans() ([]models.IPAddress, error) {
	result, err := c.datastore.GetQueries().GetIPAddressBans(context.Background())
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
func (c *ChatRepository) RemoveIPAddressBan(address string) error {
	return c.datastore.GetQueries().RemoveIPAddressBan(context.Background(), address)
}
