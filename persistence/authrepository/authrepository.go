package authrepository

import (
	"database/sql"

	"github.com/owncast/owncast/models"
)

type AuthRepository interface {
	CreateBanIPTable(db *sql.DB)
	BanIPAddress(address, note string) error
	IsIPAddressBanned(address string) (bool, error)
	GetIPAddressBans() ([]models.IPAddress, error)
	RemoveIPAddressBan(address string) error
}
