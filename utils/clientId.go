package utils

import (
	"net/http"
	"strings"
)

//GenerateClientIDFromRequest generates a client id from the provided request
func GenerateClientIDFromRequest(req *http.Request) string {
	var clientID string

	xForwardedFor := req.Header.Get("X-FORWARDED-FOR")
	if xForwardedFor != "" {
		clientID = xForwardedFor
	} else {
		ipAddressString := req.RemoteAddr
		ipAddressComponents := strings.Split(ipAddressString, ":")
		ipAddressComponents[len(ipAddressComponents)-1] = ""
		clientID = strings.Join(ipAddressComponents, ":")
	}

	// fmt.Println("IP address determined to be", ipAddress)

	return clientID + req.UserAgent()
}
