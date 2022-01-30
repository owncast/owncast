package mailjet

import (
	mj "github.com/mailjet/mailjet-apiv3-go"
	"github.com/mailjet/mailjet-apiv3-go/resources"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"

	"github.com/pkg/errors"
)

// MailJet represents an instance of the MailJet email service.
type MailJet struct {
	APIKey    string
	APISecret string
	Client    *mj.Client
}

// New returns a new instance of the MailJet email service.
func New(apiKey, apiSecret string) *MailJet {
	return &MailJet{
		APIKey:    apiKey,
		APISecret: apiSecret,
		Client:    mj.NewMailjetClient(apiKey, apiSecret),
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
	request := &mj.Request{
		Resource: "contactslist",
		ID:       listID,
		Action:   "managecontact",
	}
	fullRequest := &mj.FullRequest{
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
	mr := &mj.Request{
		Resource: "contactslist",
	}
	fmr := &mj.FullRequest{
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

// AddEmailToList will add a single email address to a Mailjet email list,
// creating a list if one does not exist.
func AddEmailToList(emailAddress string) error {
	mailjetConfig := data.GetMailjetConfiguration()
	smtpConfig := data.GetSMTPConfiguration()

	m := New(smtpConfig.Username, smtpConfig.Password)

	// If we have not previously created an email list for Owncast then create
	// a new one now, and add the requested email address to it.
	if mailjetConfig.ListID == 0 {
		listAddress, listID, err := m.CreateListAndAddAddress(emailAddress)
		if err != nil {
			log.Errorln(err)
			return errors.Wrap(err, "unable to create email list")
		}
		smtpConfig.ListAddress = listAddress
		mailjetConfig.ListID = listID
		if err := data.SetMailjetConfiguration(mailjetConfig); err != nil {
			return errors.Wrap(err, "unable to save mailjet configuration")
		}
		if err := data.SetSMTPConfiguration(smtpConfig); err != nil {
			return errors.Wrap(err, "unable to save smtp configuration")
		}
	} else {
		if err := m.AddEmailToList(emailAddress, mailjetConfig.ListID); err != nil {
			return errors.Wrap(err, "unable to add email address to list")
		}
	}
	return nil
}
