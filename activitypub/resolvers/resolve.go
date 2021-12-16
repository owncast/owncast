package resolvers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Resolve will translate a raw ActivityPub payload and fire the callback associated with that activity type.
func Resolve(c context.Context, data []byte, callbacks ...interface{}) error {
	jsonResolver, err := streams.NewJSONResolver(callbacks...)
	if err != nil {
		// Something in the setup was wrong. For example, a callback has an
		// unsupported signature and would never be called
		return err
	}

	var jsonMap map[string]interface{}
	if err = json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	log.Debugln("Resolving payload...", string(data))

	// The createCallback function will be called.
	err = jsonResolver.Resolve(c, jsonMap)
	if err != nil && !streams.IsUnmatchedErr(err) {
		// Something went wrong
		return err
	} else if streams.IsUnmatchedErr(err) {
		// Everything went right but the callback didn't match or the ActivityStreams
		// type is one that wasn't code generated.
		log.Debugln("No match: ", err)
	}

	return nil
}

// ResolveIRI will resolve an IRI ahd call the correct callback for the resolved type.
func ResolveIRI(c context.Context, iri string, callbacks ...interface{}) error {
	log.Debugln("Resolving", iri)

	req, _ := http.NewRequest("GET", iri, nil)

	actor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	if err := crypto.SignRequest(req, nil, actor); err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// fmt.Println(string(data))
	return Resolve(c, data, callbacks...)
}

// GetResolvedActorFromActorProperty resolve an actor property to a fully populated person.
func GetResolvedActorFromActorProperty(actor vocab.ActivityStreamsActorProperty) (apmodels.ActivityPubActor, error) {
	var err error
	var apActor apmodels.ActivityPubActor

	personCallback := func(c context.Context, person vocab.ActivityStreamsPerson) error {
		apActor = apmodels.MakeActorFromPerson(person)
		return nil
	}

	serviceCallback := func(c context.Context, s vocab.ActivityStreamsService) error {
		apActor = apmodels.MakeActorFromService(s)
		return nil
	}

	for iter := actor.Begin(); iter != actor.End(); iter = iter.Next() {
		if iter.IsIRI() {
			iri := iter.GetIRI()
			c := context.TODO()
			if e := ResolveIRI(c, iri.String(), personCallback, serviceCallback); e != nil {
				err = e
			}
		} else if iter.IsActivityStreamsPerson() {
			person := iter.GetActivityStreamsPerson()
			apActor = apmodels.MakeActorFromPerson(person)
		}
	}

	return apActor, errors.Wrap(err, "unable to resolve actor from actor property")
}

// GetResolvedActorFromIRI will resolve an IRI string to a fully populated actor.
func GetResolvedActorFromIRI(personOrServiceIRI string) (apmodels.ActivityPubActor, error) {
	var err error
	var apActor apmodels.ActivityPubActor

	personCallback := func(c context.Context, person vocab.ActivityStreamsPerson) error {
		apActor = apmodels.MakeActorFromPerson(person)
		return nil
	}

	serviceCallback := func(c context.Context, s vocab.ActivityStreamsService) error {
		apActor = apmodels.MakeActorFromService(s)
		return nil
	}

	if e := ResolveIRI(context.TODO(), personOrServiceIRI, personCallback, serviceCallback); e != nil {
		err = e
	}

	return apActor, errors.Wrap(err, "unable to resolve actor from IRI string")
}
