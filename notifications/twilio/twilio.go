package twilio

import (
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"
	twilioservice "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// Twilio is an instance of the Twilio notification service.
type Twilio struct {
	fromPhoneNumber string
	client          *twilio.RestClient
}

// New will create a new instance of the Twilio service given the credentials.
func New(fromPhoneNumber, accountSid, accessToken string) (*Twilio, error) {
	client := twilioservice.NewRestClientWithParams(twilioservice.RestClientParams{
		Username: accountSid,
		Password: accessToken,
	})
	return &Twilio{fromPhoneNumber: fromPhoneNumber, client: client}, nil
}

func (t *Twilio) Send(content, to string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(t.fromPhoneNumber)
	params.SetBody(content)

	_, err := t.client.ApiV2010.CreateMessage(params)
	if err != nil {
		return errors.Wrap(err, "error sending twilio message")
	}

	return nil
}
