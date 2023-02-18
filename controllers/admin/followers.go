package admin

import (
	"encoding/json"
	"net/http"
)

// ApproveFollower will approve a federated follow request.
func (c *Controller) ApproveFollower(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type approveFollowerRequest struct {
		ActorIRI string `json:"actorIRI"`
		Approved bool   `json:"approved"`
	}

	decoder := json.NewDecoder(r.Body)
	var approval approveFollowerRequest
	if err := decoder.Decode(&approval); err != nil {
		c.Service.WriteSimpleResponse(w, false, "unable to handle follower state with provided values")
		return
	}

	if approval.Approved {
		// Approve a follower
		if err := c.Follower.ApprovePreviousFollowRequest(approval.ActorIRI); err != nil {
			c.Service.WriteSimpleResponse(w, false, err.Error())
			return
		}

		localAccountName := c.Data.GetDefaultFederationUsername()

		followRequest, err := c.Follower.GetFollower(approval.ActorIRI)
		if err != nil {
			c.Service.WriteSimpleResponse(w, false, err.Error())
			return
		}

		// Send the approval to the follow requestor.
		if err := c.Follower.SendFollowAccept(followRequest.Inbox, followRequest.RequestObject, localAccountName); err != nil {
			c.Service.WriteSimpleResponse(w, false, err.Error())
			return
		}
	} else {
		// Remove/block a follower
		if err := c.Follower.BlockOrRejectFollower(approval.ActorIRI); err != nil {
			c.Service.WriteSimpleResponse(w, false, err.Error())
			return
		}
	}

	c.Service.WriteSimpleResponse(w, true, "follower updated")
}

// GetPendingFollowRequests will return a list of pending follow requests.
func (c *Controller) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := c.Follower.GetPendingFollowRequests()
	if err != nil {
		c.Service.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.Service.WriteResponse(w, requests)
}

// GetBlockedAndRejectedFollowers will return blocked and rejected followers.
func (c *Controller) GetBlockedAndRejectedFollowers(w http.ResponseWriter, r *http.Request) {
	rejections, err := c.Follower.GetBlockedAndRejectedFollowers()
	if err != nil {
		c.Service.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.Service.WriteResponse(w, rejections)
}
