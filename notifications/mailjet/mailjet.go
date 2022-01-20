package mailjet

import (
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/mailjet/mailjet-apiv3-go/resources"
	"github.com/teris-io/shortid"

	"github.com/pkg/errors"
)

// MailJet represents an instance of the MailJet email service.
type MailJet struct {
	APIKey    string
	APISecret string
	Client    *mailjet.Client
}

// New returns a new instance of the MailJet email service.
func New(apiKey, apiSecret string) *MailJet {
	return &MailJet{
		APIKey:    apiKey,
		APISecret: apiSecret,
		Client:    mailjet.NewMailjetClient(apiKey, apiSecret),
	}
}

// CreateListAndAddAddress will create a new email list and add address to it.
func (m *MailJet) CreateListAndAddAddress(address string) (string, int64, error) {
	listAddress, listID, err := m.createList("owncast-" + shortid.MustGenerate())
	if err != nil {
		return "", 0, err
	}

	if err := m.AddEmailToList(address, listID); err != nil {
		return "", 0, err
	}

	return listAddress, listID, nil
}

// AddEmailToList will add an email address to the provided list.
func (m *MailJet) AddEmailToList(address string, listID int64) error {
	var data []resources.ContactslistManageContact
	request := &mailjet.Request{
		Resource: "contactslist",
		ID:       listID,
		Action:   "managecontact",
	}
	fullRequest := &mailjet.FullRequest{
		Info: request,
		Payload: &resources.ContactslistManageContact{
			Properties: "object",
			Action:     "addnoforce",
			Email:      address,
		},
	}
	if err := m.Client.Post(fullRequest, &data); err != nil {
		return errors.Wrap(err, "unable to subscribe email to list")
	}

	return nil
}

// createList will create a new email list on Mailjet.
func (m *MailJet) createList(name string) (string, int64, error) {
	var data []resources.Contactslist
	mr := &mailjet.Request{
		Resource: "contactslist",
	}
	fmr := &mailjet.FullRequest{
		Info: mr,
		Payload: &resources.Contactslist{
			Name: name,
		},
	}
	err := m.Client.Post(fmr, &data)
	if err != nil {
		return "", 0, errors.Wrap(err, "unable to create email list from provider")
	}

	if len(data) == 0 {
		return "", 0, errors.New("provider returned no new email lists")
	}

	list := data[0]
	return list.Address + "@lists.mailjet.com", list.ID, nil
}
