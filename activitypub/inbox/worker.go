package inbox

import (
	"context"
	"fmt"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/resolvers"
)

var ch = make(chan []byte, 5)

// func Run() {
// 	for v := range ch {
// 		fmt.Println("run...")
// 		h := handle(v)
// 		handled := <-h
// 		fmt.Println("Handled", handled)
// 		// Save followRequest
// 		// Send ACCEPT back to actor
// 	}
// }

func Add(data []byte, forLocalAccount string) {
	fmt.Println("Adding AP Payload...")
	handle(data, forLocalAccount)
	// ch <- data
}

func handle(data []byte, forLocalAccount string) chan bool {
	c := context.WithValue(context.Background(), "account", forLocalAccount)
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
		// fmt.Println(activity)
		r <- false
		return nil
	}

	// go func() {
	if err := resolvers.Resolve(data, c, createCallback, updateCallback, handleFollowInboxRequest, personCallback, handleUndoInboxRequest); err != nil {
		panic(err)
	}
	// }()

	return r
}
