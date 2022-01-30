package email

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
)

//go:embed "golive.tmpl.html"
var goLiveTemplate string

// Email represents an instance of the Email notifier.
type Email struct {
	From       string
	SMTPServer string
	SMTPPort   string
	Username   string
	Password   string
}

// New creates a new instance of the Email notifier.
func new(from, server, port, username, password string) *Email {
	return &Email{
		From:       from,
		SMTPServer: server,
		SMTPPort:   port,
		Username:   username,
		Password:   password,
	}
}

// New creates a new instance of the email notifier.
func New() (*Email, error) {
	smtpConfig := data.GetSMTPConfiguration()
	if smtpConfig.Enabled && smtpConfig.FromAddress != "" {
		e := new(smtpConfig.FromAddress, smtpConfig.Server, "587", smtpConfig.Username, smtpConfig.Password)
		return e, nil
	}

	return nil, errors.New("email delivery not configured")
}

// Send will send an email notification.
func (e *Email) Send(to []string, content, subject string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", e.From)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(to, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s\r\n", content)

	return smtp.SendMail(e.SMTPServer+":"+e.SMTPPort, auth, e.From, to, []byte(msg))
}

// GenerateEmailContent will return email content as a string.
func GenerateEmailContent() (string, error) {
	type templateData struct {
		Logo              string
		Thumbnail         string
		ServerURL         string
		ServerName        string
		Description       string
		StreamDescription string
	}

	td := templateData{
		Logo:              data.GetServerURL() + "/logo",
		Thumbnail:         data.GetServerURL() + "/thumbnail.jpg",
		ServerURL:         data.GetServerURL(),
		ServerName:        data.GetServerName(),
		Description:       data.GetServerSummary(),
		StreamDescription: data.GetStreamTitle(),
	}

	t, err := template.New("goLive").Parse(goLiveTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse go live email template")
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, td); err != nil {
		return "", errors.Wrap(err, "failed to execute go live email template")
	}

	content := tpl.String()

	return content, nil
}
