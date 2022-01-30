package email

import (
	"github.com/owncast/owncast/static"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

// Email addresses with a double opt-in confirmation pending. This is stored
// in memory only, as to not store an email address in the database.
var pending = map[string]string{}

// CreateDoubleOptInRequest will add a pending double opt-in request.
func CreateDoubleOptInRequest(email string) (string, error) {
	// Create a unique ID for the request.
	id, err := shortid.Generate()
	if err != nil {
		return "", err
	}

	// Store the email address and ID.
	pending[id] = email

	// Send confirmation email.
	link := "http://localhost:8080/email/confirm?id=" + id
	if err := sendConfirmationEmail(email, link); err != nil {
		return "", errors.Wrap(err, "unable to send confirmation email")
	}

	return id, nil
}

// HandleRequest will handle a double opt-in request.
func HandleRequest(id string) (string, error) {
	// Get the email address.
	email, ok := pending[id]
	if !ok {
		return email, errors.New("Email confirmation request not found")
	}

	// Delete the request.
	delete(pending, id)

	return email, nil
}

func sendConfirmationEmail(emailAddress, link string) error {
	content, err := static.GetEmailConfirmTemplateWithContent("Confirm your email address", "Please confirm your email address", "", "Confirm", link)
	if err != nil {
		return errors.Wrap(err, "unable to generate email content")
	}

	e, err := New()
	if err != nil {
		return errors.Wrap(err, "unable to create email service")
	}

	return e.Send([]string{emailAddress}, content, "Confirm your email address")
}
