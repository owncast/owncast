// This package utilizes the MaxMind GeoLite2 GeoIP database https://dev.maxmind.com/geoip/geoip2/geolite2/.
// You must provide your own copy of this database for it to work.
// Read more about how this works at http://owncast.online/docs/geoip

package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

var _geoIPCache = map[string]GeoDetails{}
var _enabled = true // Try to use GeoIP support it by default.
var geoIPDatabasePath = "data/GeoLite2-City.mmdb"

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

	if ip == "::1" || ip == "127.0.0.1" {
		return &GeoDetails{
			CountryCode: "N/A",
			RegionName:  "Localhost",
			TimeZone:    "",
		}
	}

	return fetchGeoForIP(ip)
}

// fetchGeoForIP makes an API call to get geo details for an IP address.
func fetchGeoForIP(ip string) *GeoDetails {
	// If GeoIP has been disabled then don't try to access it.
	if !_enabled {
		return nil
	}

	// Don't re-fetch if we already have it.
	if geoDetails, ok := _geoIPCache[ip]; ok {
		return &geoDetails
	}

	db, err := geoip2.Open(geoIPDatabasePath)
	if err != nil {
		log.Traceln("GeoIP support is disabled. visit http://owncast.online/docs/geoip to learn how to enable.", err)
		_enabled = false
		return nil
	}

	defer db.Close()

	ipObject := net.ParseIP(ip)

	record, err := db.City(ipObject)
	if err != nil {
		log.Warnln(err)
		return nil
	}

	// If no country is available then exit
	if record.Country.IsoCode == "" {
		return nil
	}

	// If we believe this IP to be anonymous then no reason to report it
	if record.Traits.IsAnonymousProxy {
		return nil
	}

	var regionName = "Unknown"
	if len(record.Subdivisions) > 0 {
		if region, ok := record.Subdivisions[0].Names["en"]; ok {
			regionName = region
		}
	}

	response := GeoDetails{
		CountryCode: record.Country.IsoCode,
		RegionName:  regionName,
		TimeZone:    record.Location.TimeZone,
	}

	_geoIPCache[ip] = response

	return &response
}
