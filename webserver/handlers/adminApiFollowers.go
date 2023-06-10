package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/activitypub/persistence"
	aprequests "github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"
)

// ApproveFollower will approve a federated follow request.
func (h *Handlers) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type approveFollowerRequest struct {
		ActorIRI string `json:"actorIRI"`
		Approved bool   `json:"approved"`
	}

	decoder := json.NewDecoder(r.Body)
	var approval approveFollowerRequest
	if err := decoder.Decode(&approval); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to handle follower state with provided values")
		return
	}

	if approval.Approved {
		// Approve a follower
		if err := persistence.ApprovePreviousFollowRequest(approval.ActorIRI); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		localAccountName := data.GetDefaultFederationUsername()

		followRequest, err := persistence.GetFollower(approval.ActorIRI)
		if err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		// Send the approval to the follow requestor.
		if err := aprequests.SendFollowAccept(followRequest.Inbox, followRequest.RequestObject, localAccountName); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}
	} else {
		// Remove/block a follower
		if err := persistence.BlockOrRejectFollower(approval.ActorIRI); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	responses.WriteSimpleResponse(w, true, "follower updated")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func (h *Handlers) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := persistence.GetPendingFollowRequests()
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteResponse(w, requests)
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func (h *Handlers) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	rejections, err := persistence.GetBlockedAndRejectedFollowers()
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteResponse(w, rejections)
}
