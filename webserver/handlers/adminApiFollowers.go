package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/services/apfederation/outbox"
	"github.com/owncast/owncast/storage/federationrepository"
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

	federationRepository := federationrepository.Get()

	if approval.Approved {
		// Approve a follower
		if err := federationRepository.ApprovePreviousFollowRequest(approval.ActorIRI); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		localAccountName := configRepository.GetDefaultFederationUsername()

		followRequest, err := federationRepository.GetFollower(approval.ActorIRI)
		if err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}

		ob := outbox.Get()

		// Send the approval to the follow requestor.
		if err := ob.SendFollowAccept(followRequest.Inbox, followRequest.RequestObject, localAccountName); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}
	} else {
		// Remove/block a follower
		if err := federationRepository.BlockOrRejectFollower(approval.ActorIRI); err != nil {
			responses.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	responses.WriteSimpleResponse(w, true, "follower updated")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func (h *Handlers) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	federationRepository := federationrepository.Get()
	requests, err := federationRepository.GetPendingFollowRequests()
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteResponse(w, requests)
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func (h *Handlers) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	federationRepository := federationrepository.Get()
	rejections, err := federationRepository.GetBlockedAndRejectedFollowers()
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteResponse(w, rejections)
}
