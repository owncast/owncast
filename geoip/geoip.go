// This package utilizes the MaxMind GeoLite2 GeoIP database https://dev.maxmind.com/geoip/geoip2/geolite2/.
// You must provide your own copy of this database for it to work.
// Read more about how this works at http://owncast.online/docs/geoip

package geoip

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

const geoIPDatabasePath = "data/GeoLite2-City.mmdb"

// Client can look up geography information for IP addresses.
type Client struct {
	cache   sync.Map
	enabled int32
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		enabled: 1, // Try to use GeoIP support by default.
	}
}

// GeoDetails stores details about a location.
type GeoDetails struct {
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	TimeZone    string `json:"timeZone"`
}

// GetGeoFromIP returns geo details associated with an IP address if we
// have previously fetched it.
func (c *Client) GetGeoFromIP(ip string) *GeoDetails {
	if cachedGeoDetails, ok := c.cache.Load(ip); ok {
		return cachedGeoDetails.(*GeoDetails)
	}

	if ip == "::1" || ip == "127.0.0.1" {
		return &GeoDetails{
			CountryCode: "N/A",
			RegionName:  "Localhost",
			TimeZone:    "",
		}
	}

	return c.fetchGeoForIP(ip)
}

// fetchGeoForIP makes an API call to get geo details for an IP address.
func (c *Client) fetchGeoForIP(ip string) *GeoDetails {
	// If GeoIP has been disabled then don't try to access it.
	if atomic.LoadInt32(&c.enabled) == 0 {
		return nil
	}

	db, err := geoip2.Open(geoIPDatabasePath)
	if err != nil {
		log.Traceln("GeoIP support is disabled. visit https://owncast.online/docs/geoip to learn how to enable.", err)
		atomic.StoreInt32(&c.enabled, 0)
		return nil
	}
	defer db.Close()

	var response *GeoDetails
	ipObject := net.ParseIP(ip)

	record, err := db.City(ipObject)
	if err == nil {
		// If no country is available then exit
		// If we believe this IP to be anonymous then no reason to report it
		if record.Country.IsoCode != "" && !record.Traits.IsAnonymousProxy {
			var regionName = "Unknown"
			if len(record.Subdivisions) > 0 {
				if region, ok := record.Subdivisions[0].Names["en"]; ok {
					regionName = region
				}
			}

			response = &GeoDetails{
				CountryCode: record.Country.IsoCode,
				RegionName:  regionName,
				TimeZone:    record.Location.TimeZone,
			}
		}
	} else {
		log.Warnln(err)
	}

	c.cache.Store(ip, response)

	return response
}
