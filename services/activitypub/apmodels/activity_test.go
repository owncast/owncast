package apmodels

import (
	"net/url"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
)

func TestDirectMessageAddressing(t *testing.T) {
	recipient, err := url.Parse("https://remote.example/users/test")
	if err != nil {
		t.Fatal(err)
	}

	note := MakeNoteDirect(streams.NewActivityStreamsNote(), recipient)
	if to := note.GetActivityStreamsTo(); to == nil || to.Len() != 1 || to.At(0).GetIRI().String() != recipient.String() {
		t.Fatalf("note recipient = %v, want %s", to, recipient)
	}
	if note.GetActivityStreamsCc() != nil {
		t.Fatal("note direct recipient should not be cc'd")
	}

	activity := MakeActivityDirect(streams.NewActivityStreamsCreate(), recipient)
	if to := activity.GetActivityStreamsTo(); to == nil || to.Len() != 1 || to.At(0).GetIRI().String() != recipient.String() {
		t.Fatalf("activity recipient = %v, want %s", to, recipient)
	}
	if activity.GetActivityStreamsCc() != nil {
		t.Fatal("activity direct recipient should not be cc'd")
	}
}
