package user

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

// This stores frequently accessed data about chat users.
var _userIdCache = ttlcache.NewCache()
var _userAccessTokenCache = ttlcache.NewCache()

func init() {
	_userIdCache.SetTTL(time.Duration(10 * time.Minute))
	_userAccessTokenCache.SetTTL(time.Duration(10 * time.Minute))
}

func invalidateUserCache(id string) {
	_userIdCache.Remove(id)
	_userAccessTokenCache.Remove(id)
}

func getCachedIdUser(id string) *User {
	if entry, error := _userIdCache.Get(id); error == nil {
		return entry.(*User)
	}
	return nil
}

func setCachedIdUser(id string, user *User) {
	_userIdCache.Set(id, user)
}

func getCachedAccessTokenUser(id string) *User {
	if entry, err := _userAccessTokenCache.Get(id); err == nil {
		return entry.(*User)
	}
	return nil
}

func setCachedAccessTokenUser(id string, user *User) {
	_userAccessTokenCache.Set(id, user)
}
