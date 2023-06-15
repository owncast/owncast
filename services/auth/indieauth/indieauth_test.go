package indieauth

import (
	"testing"

	"github.com/owncast/owncast/utils"
)

var indieAuthServer = GetIndieAuthServer()

func TestLimitGlobalPendingRequests(t *testing.T) {
	// Simulate 10 pending requests
	for i := 0; i < maxPendingRequests-1; i++ {
		cid, _ := utils.GenerateRandomString(10)
		redirectURL, _ := utils.GenerateRandomString(10)
		cc, _ := utils.GenerateRandomString(10)
		state, _ := utils.GenerateRandomString(10)
		me, _ := utils.GenerateRandomString(10)

		_, err := indieAuthServer.StartServerAuth(cid, redirectURL, cc, state, me)
		if err != nil {
			t.Error("Registration should be permitted.", i, " of ", len(indieAuthServer.pendingServerAuthRequests), err)
		}
	}

	// This should throw an error
	cid, _ := utils.GenerateRandomString(10)
	redirectURL, _ := utils.GenerateRandomString(10)
	cc, _ := utils.GenerateRandomString(10)
	state, _ := utils.GenerateRandomString(10)
	me, _ := utils.GenerateRandomString(10)

	_, err := indieAuthServer.StartServerAuth(cid, redirectURL, cc, state, me)
	if err == nil {
		t.Error("Registration should not be permitted.")
	}
}
