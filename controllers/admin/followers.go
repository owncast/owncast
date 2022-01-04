package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
)

// ApproveFollower will approve a federated follow request.
func ApproveFollower(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type approveFollowerRequest struct {
		ActorIRI string `json:"actorIRI"`
		Approved bool   `json:"approved"`
	}

	decoder := json.NewDecoder(r.Body)
	var approval approveFollowerRequest
	if err := decoder.Decode(&approval); err != nil {
		controllers.WriteSimpleResponse(w, false, "unable to handle follower state with provided values")
		return
	}

	if approval.Approved {
		// Approve a follower
		if err := persistence.ApprovePreviousFollowRequest(approval.ActorIRI); err != nil {
			controllers.WriteSimpleResponse(w, false, err.Error())
			return
		}

		localAccountName := data.GetDefaultFederationUsername()

		follower, err := persistence.GetFollower(approval.ActorIRI)
		if err != nil {
			controllers.WriteSimpleResponse(w, false, err.Error())
			return
		}

		// Send the approval to the follow requestor.
		if err := requests.SendFollowAccept(follower.Inbox, follower.FollowRequestIri, localAccountName); err != nil {
			controllers.WriteSimpleResponse(w, false, err.Error())
			return
		}
	} else {
		// Remove/block a follower
		if err := persistence.BlockOrRejectFollower(approval.ActorIRI); err != nil {
			controllers.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	controllers.WriteSimpleResponse(w, true, "follower updated")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := persistence.GetPendingFollowRequests()
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteResponse(w, requests)
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	rejections, err := persistence.GetBlockedAndRejectedFollowers()
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteResponse(w, rejections)
}
