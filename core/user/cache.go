package user

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	log "github.com/sirupsen/logrus"
)

// This stores frequently accessed data about chat users.
var _userIdCache = ttlcache.NewCache()
var _userAccessTokenCache = ttlcache.NewCache()

func init() {
	if err := _userIdCache.SetTTL(10 * time.Minute); err != nil {
		log.Debugln(err)
	}
	if err := _userAccessTokenCache.SetTTL(10 * time.Minute); err != nil {
		log.Debugln(err)
	}
}

func invalidateUserCache(id string) {
	if err := _userIdCache.Remove(id); err != nil {
		log.Debugln(err)
	}
	if err := _userAccessTokenCache.Remove(id); err != nil {
		log.Debugln(err)
	}
}

func getCachedIdUser(id string) *User {
	if entry, error := _userIdCache.Get(id); error == nil {
		return entry.(*User)
	}
	return nil
}

func setCachedIdUser(id string, user *User) {
	if err := _userIdCache.Set(id, user); err != nil {
		log.Debugln(err)
	}
}

func getCachedAccessTokenUser(id string) *User {
	if entry, err := _userAccessTokenCache.Get(id); err == nil {
		return entry.(*User)
	}
	return nil
}

func setCachedAccessTokenUser(id string, user *User) {
	if err := _userAccessTokenCache.Set(id, user); err != nil {
		log.Debugln(err)
	}
}
