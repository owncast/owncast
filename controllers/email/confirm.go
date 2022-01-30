package email

import (
	"fmt"
	"net/http"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/notifications/email"
	"github.com/owncast/owncast/notifications/mailjet"

	"github.com/owncast/owncast/static"
	log "github.com/sirupsen/logrus"
)

// ConfirmEmailDoubleOptIn is accessed in response to a verification email
// that was sent to an email address confirming they want to to be notified
// by email.
func ConfirmEmailDoubleOptIn(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	emailAddress, err := email.HandleRequest(id)

	var title string
	var subtitle string
	var footer string

	if err == nil {
		title = fmt.Sprintf("We will let you know when %s goes live", data.GetServerName())
		footer = "By the way, we don't store your email. Only our email provider has access to it."

		// Add the confirmed email to the mailining list.
		if err := mailjet.AddEmailToList(emailAddress); err != nil {
			title = "There was a problem."
			subtitle = "Unable to add your email to the notification list."
			footer = err.Error()
			log.Errorln(err)
		}
	} else if emailAddress == "" || err != nil {
		title = "Could not confirm your email address."
		footer = "You can try again to get notifications."
	}

	// Return the confirmation page.
	content, err := static.GetEmailConfirmTemplateWithContent(title, subtitle, footer, "", "")
	if err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(content)); err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
