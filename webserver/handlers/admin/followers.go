package admin

import (
	"encoding/json"
	"net/http"

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
	if approval.ActorIRI == nil || approval.Approved == nil {
		webutils.WriteSimpleResponse(w, false, "actorIRI and approved are required")
		return
	}

	streamActive := a.stream.GetStatus().Online
	followRequest, err := a.activitypub.RespondToFollow(*approval.ActorIRI, *approval.Approved, streamActive)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	if *approval.Approved && followRequest.ApprovedAt == nil && !followRequest.IsDirectory {
		go a.webhooks.SendFediverseEngagementFollowEvent(*approval.ActorIRI)
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

	if err := a.activitypub.RemoveFollower(request.ActorIRI); err != nil {
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

// GetDirectoryFollowers returns the directories that are featuring/listing this
// server: approved followers that identified themselves with the ns#directory
// marker. The operator can review them here and remove any with the existing
// RemoveFollower endpoint.
func (a *Admin) GetDirectoryFollowers(w http.ResponseWriter, r *http.Request) {
	followers, err := a.followersRepository.GetApprovedDirectoryFollowers()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteResponse(w, followers)
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
