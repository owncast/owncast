package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/core/data"
)

// Resolve will translate a raw ActivityPub payload and fire the callback associated with that activity type.
func Resolve(data []byte, c context.Context, callbacks ...interface{}) error {
	fmt.Println("Resolving payload...")
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

	// fmt.Println(string(data))

	// The createCallback function will be called.
	err = jsonResolver.Resolve(c, jsonMap)
	if err != nil && !streams.IsUnmatchedErr(err) {
		// Something went wrong
		return err
	} else if streams.IsUnmatchedErr(err) {
		// Everything went right but the callback didn't match or the ActivityStreams
		// type is one that wasn't code generated.
		fmt.Println("No match: ", err)
	}

	return nil
}

func ResolveIRI(iri string, c context.Context, callbacks ...interface{}) error {
	fmt.Println("Resolving", iri)

	req, _ := http.NewRequest("GET", iri, nil)

	actor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	if err := crypto.SignRequest(req, nil, actor); err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// fmt.Println(string(data))
	return Resolve(data, c, callbacks...)
}

func GetResolvedPersonFromActor(actor vocab.ActivityStreamsActorProperty) (vocab.ActivityStreamsPerson, error) {
	var err error
	var person vocab.ActivityStreamsPerson

	personCallback := func(c context.Context, p vocab.ActivityStreamsPerson) error {
		person = p
		return nil
	}

	for iter := actor.Begin(); iter != actor.End(); iter = iter.Next() {
		if iter.IsIRI() {
			fmt.Println("Is IRI", iter.GetIRI())
			iri := iter.GetIRI()
			c := context.TODO()
			if e := ResolveIRI(iri.String(), c, personCallback); e != nil {
				err = e
			}
		} else if iter.IsActivityStreamsPerson() {
			p := iter.GetActivityStreamsPerson()
			person = p
		}
	}

	return person, err
}
