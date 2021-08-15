package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/resolvers"
)

var _queue = make(chan models.InboxRequest, 5)

func init() {
	go run()
}

func run() {
	for r := range _queue {
		fmt.Println("run...")
		handle(r)
	}
}

func Add(request models.InboxRequest) {
	fmt.Println("Adding AP Payload...")
	_queue <- request
}

func handle(request models.InboxRequest) chan bool {
	c := context.WithValue(context.Background(), "account", request.ForLocalAccount)
	r := make(chan bool)

	fmt.Println("Handling payload via worker...")

	createCallback := func(c context.Context, activity vocab.ActivityStreamsCreate) error {
		fmt.Println("createCallback fired!")

		fmt.Println(activity)
		r <- false
		return nil
	}

	updateCallback := func(c context.Context, activity vocab.ActivityStreamsUpdate) error {
		fmt.Println("updateCallback fired!")

		fmt.Println(activity)
		r <- false
		return nil
	}

	personCallback := func(c context.Context, activity vocab.ActivityStreamsPerson) error {
		fmt.Println("personCallback fired!")
		r <- false
		return nil
	}

	if err := resolvers.Resolve(request.Data, c, createCallback, updateCallback, handleFollowInboxRequest, personCallback, handleUndoInboxRequest); err != nil {
		panic(err)
	}

	return r
}
