package mailer

import (
	"bytes"
	"errors"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(FromEmail, APIKey string) (Mailer, error) {
	APIKey = os.Getenv("API_KEY")
	if APIKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: FromEmail,
		apiKey:    APIKey,
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
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)
	return dialer.DialAndSend(message)
}
