package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/services/activitypub/requests"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ApproveFollower will approve a federated follow request.
func (a *Admin) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	// type approveFollowerRequest struct {
	// 	ActorIRI string `json:"actorIRI"`
	// 	Approved bool   `json:"approved"`
	// }

	decoder := json.NewDecoder(r.Body)
	var approval generated.ApproveFollowerJSONBody
	if err := decoder.Decode(&approval); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to handle follower state with provided values")
		return
	}

	if *approval.Approved {
		// Approve a follower
		if err := a.followersRepository.ApprovePreviousRequest(*approval.ActorIRI); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		localAccountName := a.configRepository.GetDefaultFederationUsername()

		followRequest, err := a.followersRepository.GetByIRI(*approval.ActorIRI)
		if err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		// Featured-streams follows are a directory relationship, not a fan
		// follow, so don't fire the follower webhook for them.
		if !followRequest.IsOwncastServer {
			go a.webhooks.SendFediverseEngagementFollowEvent(*approval.ActorIRI)
		}

		// Send the approval to the follow requestor, including our current
		// stream status so a featured-streams follower approved while we are
		// already live shows us live immediately.
		streamActive := a.stream.GetStatus().Online
		if err := requests.SendFollowAccept(a.activitypub.Workerpool(), followRequest.Inbox, followRequest.RequestObject, localAccountName, a.apBuilder, a.apSigner, a.configRepository, streamActive); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}
	} else {
		// Remove/block a follower
		if err := a.followersRepository.BlockOrReject(*approval.ActorIRI); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	webutils.WriteSimpleResponse(w, true, "follower updated")
}

// RemoveFollower removes a follower without blocking them. Unlike rejecting,
// this deletes the follow outright (no disabled_at), so the actor is free to
// follow again later.
func (a *Admin) RemoveFollower(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	var request generated.RemoveFollowerJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to parse request: "+err.Error())
		return
	}
	if request.ActorIRI == "" {
		webutils.WriteSimpleResponse(w, false, "actorIRI is required")
		return
	}

	if err := a.followersRepository.RemoveByIRI(request.ActorIRI); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "follower removed")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func (a *Admin) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := a.followersRepository.GetPendingFollowRequests()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteResponse(w, requests)
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func (a *Admin) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	rejections, err := a.followersRepository.GetBlockedAndRejected()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteResponse(w, rejections)
}
