package indieauth

import (
	"testing"

	"github.com/owncast/owncast/utils"
)

func TestLimitGlobalPendingRequests(t *testing.T) {
	// Construct an isolated Service for this test. CompleteServerAuth is
	// not called here, so the ConfigRepository can be nil.
	svc := New(Deps{ConfigRepository: nil})

	// Simulate maxPendingRequests-1 pending requests.
	for i := 0; i < maxPendingRequests-1; i++ {
		cid, _ := utils.GenerateRandomString(10)
		redirectURL, _ := utils.GenerateRandomString(10)
		cc, _ := utils.GenerateRandomString(10)
		state, _ := utils.GenerateRandomString(10)
		me, _ := utils.GenerateRandomString(10)

		_, err := svc.StartServerAuth(cid, redirectURL, cc, state, me)
		if err != nil {
			t.Error("Registration should be permitted.", i, " of ", len(svc.pendingServerAuthRequests), err)
		}
	}

	// This should throw an error.
	cid, _ := utils.GenerateRandomString(10)
	redirectURL, _ := utils.GenerateRandomString(10)
	cc, _ := utils.GenerateRandomString(10)
	state, _ := utils.GenerateRandomString(10)
	me, _ := utils.GenerateRandomString(10)

	_, err := svc.StartServerAuth(cid, redirectURL, cc, state, me)
	if err == nil {
		t.Error("Registration should not be permitted.")
	}
}
