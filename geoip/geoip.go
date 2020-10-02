package geoip

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"
)

var geoIPCache = map[string]GeoDetails{}

type geoRequestResponse struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
	TimeZone    string `json:"time_zone"`
}

type GeoDetails struct {
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	RegionCode  string `json:"regionCode"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	TimeZone    string `json:"timeZone"`
}

// GetGeoFromIP returns geo details associated with an IP address if we
// have previously fetched it.
func GetGeoFromIP(ip string) *GeoDetails {
	if cachedGeoDetails, ok := geoIPCache[ip]; ok {
		return &cachedGeoDetails
	}

	return nil
}

// FetchGeoForIP makes an API call to get geo details for an IP address.
func FetchGeoForIP(ip string) {
	// Don't re-fetch if we already have it.
	if _, ok := geoIPCache[ip]; ok {
		return
	}

	url := "https://freegeoip.app/json/"

	go func() {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var apiResponse geoRequestResponse
		err := json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Errorln(err)
			return
		}

		response := GeoDetails(apiResponse)
		geoIPCache[ip] = response
	}()

}
