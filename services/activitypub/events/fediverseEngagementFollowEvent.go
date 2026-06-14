package events

type FediverseEngagementFollowEvent struct {
	Event
	Name     string `json:"name"`
	Username string `json:"username"`
	Image    string `json:"image"`
}
