// This package utilizes the MaxMind GeoLite2 GeoIP database https://dev.maxmind.com/geoip/geoip2/geolite2/.
// You must provide your own copy of this database for it to work.
// Read more about how this works at http://owncast.online/docs/geoip

package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/owncast/owncast/config"
	log "github.com/sirupsen/logrus"
)

var _geoIPCache = map[string]GeoDetails{}
var _enabled = true // Try to use GeoIP support it by default.

// GeoDetails stores details about a location.
type GeoDetails struct {
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	TimeZone    string `json:"timeZone"`
}

// GetGeoFromIP returns geo details associated with an IP address if we
// have previously fetched it.
func GetGeoFromIP(ip string) *GeoDetails {
	if cachedGeoDetails, ok := _geoIPCache[ip]; ok {
		return &cachedGeoDetails
	}

	return nil
}

// FetchGeoForIP makes an API call to get geo details for an IP address.
func FetchGeoForIP(ip string) {
	// If GeoIP has been disabled then don't try to access it.
	if !_enabled {
		return
	}

	// Don't re-fetch if we already have it.
	if _, ok := _geoIPCache[ip]; ok {
		return
	}

	go func() {
		db, err := geoip2.Open(config.GeoIPDatabasePath)
		if err != nil {
			log.Traceln("GeoIP support is disabled. visit http://owncast.online/docs/geoip to learn how to enable.", err)
			_enabled = false
			return
		}

		defer db.Close()

		ipObject := net.ParseIP(ip)

		record, err := db.City(ipObject)
		if err != nil {
			log.Warnln(err)
			return
		}

		// If no country is available then exit
		if record.Country.IsoCode == "" {
			return
		}

		// If we believe this IP to be anonymous then no reason to report it
		if record.Traits.IsAnonymousProxy {
			return
		}

		response := GeoDetails{
			CountryCode: record.Country.IsoCode,
			RegionName:  record.Subdivisions[0].Names["en"],
			TimeZone:    record.Location.TimeZone,
		}

		_geoIPCache[ip] = response
	}()

}
