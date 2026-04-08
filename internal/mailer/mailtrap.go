package mailer

import (
	"bytes"
	"errors"
	"text/template"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	FromEmail string
	APIKey    string
}

func NewMailTrapClient(FromEmail, APIKey string) (Mailer, error) {
	if APIKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		FromEmail: FromEmail,
		APIKey:    APIKey,
	}, nil
}

func (m mailtrapClient) Send(username, email string, data any) error {
	tmpl, err := template.ParseFS(FS, "templates/user_invitation.tmpl")
	if err != nil {
		return err
	}

	payload := map[string]any{
		"Username":    username,
		"InvoiceLink": data,
	}

	body := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(body, "body", payload); err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(subject, "subject", payload); err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.FromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.APIKey)
	return dialer.DialAndSend(message)
}
