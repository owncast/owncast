package authrepository

import (
	"github.com/owncast/owncast/webserver/handlers/generated"
)

type AuthRepository interface {
	BanIPAddress(address, note string) error
	IsIPAddressBanned(address string) (bool, error)
	GetIPAddressBans() ([]generated.IPAddress, error)
	RemoveIPAddressBan(address string) error
}
