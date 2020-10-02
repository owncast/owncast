package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

//GenerateClientIDFromRequest generates a client id from the provided request
func GenerateClientIDFromRequest(req *http.Request) string {
	ipAddress := GetIPAddressFromRequest(req)
	ipAddressComponents := strings.Split(ipAddress, ":")
	ipAddressComponents[len(ipAddressComponents)-1] = ""
	clientID := strings.Join(ipAddressComponents, ":") + req.UserAgent()

	// Create a MD5 hash of this ip + useragent
	hasher := md5.New()
	hasher.Write([]byte(clientID))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetIPAddressFromRequest returns the IP address from a http request
func GetIPAddressFromRequest(req *http.Request) string {
	ipAddressString := req.RemoteAddr
	xForwardedFor := req.Header.Get("X-FORWARDED-FOR")
	if xForwardedFor != "" {
		ipAddressString = xForwardedFor
	}

	ip, _, err := net.SplitHostPort(ipAddressString)
	if err != nil {
		log.Errorln(err)
		return ""
	}

	return ip
}
