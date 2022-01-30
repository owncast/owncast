package static

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed "emailconfirm.html.tmpl"
var emailConfirmTemplate []byte

// GetEmailConfirmTemplate will return the email confirmation template.
func GetEmailConfirmTemplate() string {
	return string(getFileSystemStaticFileOrDefault("emailconfirm.html.tmpl", emailConfirmTemplate))
}

// GetEmailConfirmTemplateWithContent will return the email confirmation
// template with content populated.
func GetEmailConfirmTemplateWithContent(title, subtitle, footer, button, buttonLink string) (string, error) {
	type content struct {
		Title         string
		Subtitle      string
		SecondaryText string
		Button        string
		ButtonLink    string
		ShowFooter    bool
	}

	data := content{
		Title:         title,
		Subtitle:      subtitle,
		SecondaryText: footer,
		Button:        button,
		ButtonLink:    buttonLink,
		ShowFooter:    false,
	}

	templateBytes := GetEmailConfirmTemplate()

	t, err := template.New("emailConfirm").Parse(templateBytes)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
