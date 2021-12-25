package inbox

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/resolvers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

func handleEngagementActivity(eventType events.EventType, object vocab.ActivityStreamsObjectProperty, actorReference vocab.ActivityStreamsActorProperty, activityIRI *url.URL, action string) error {
	// Do nothing if displaying engagement actions has been turned off.
	if !data.GetFederationShowEngagement() {
		return nil
	}

	// Do nothing if chat is disabled
	if data.GetChatDisabled() {
		return nil
	}

	IRI := object.At(0).GetIRI().String()

	// Allow all Follows to be handled, otherwise the object IRI must match something
	// we have sent before.
	post, isLiveNotification, err := persistence.GetObjectByIRI(IRI)
	if action != "follow" && (err != nil || post == "") {
		log.Errorln("Could not find post locally:", IRI, err)
		return err
	}

	// Get actor of the Like
	actor, _ := resolvers.GetResolvedActorFromActorProperty(actorReference)

	// Send chat message
	actorName := actor.Name
	if actorName == "" {
		actorName = actor.Username
	}
	actorIRI := actorReference.Begin().GetIRI().String()

	userPrefix := fmt.Sprintf("[%s](%s) just", actorName, actorIRI)
	var suffix string
	if isLiveNotification && action == events.FediverseEngagementLike {
		suffix = fmt.Sprintf("liked that %s went live.", data.GetServerName())
	} else if isLiveNotification && action == events.FediverseEngagementRepost {
		suffix = fmt.Sprintf("shared that %s went live.", data.GetServerName())
	} else if action == events.FediverseEngagementFollow {
		suffix = fmt.Sprintf("followed %s.", data.GetServerName())
	} else {
		return fmt.Errorf("could not handle event for sending to chat: %s", action)
	}
	body := fmt.Sprintf("%s %s", userPrefix, suffix)

	var image *string
	if actor.Image != nil {
		s := actor.Image.String()
		image = &s
	}

	if err := chat.SendFediverseAction(eventType, actor.FullUsername, image, body, actorIRI); err != nil {
		return err
	}

	return nil
}
