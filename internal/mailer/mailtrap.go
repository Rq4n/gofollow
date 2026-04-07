package mailer

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(fromEmail, apiKey string) (Mailer, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

// data interface -> invoice_link
func (m mailtrapClient) Send(username, email string, data any) error {
	// tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	// if err != nil {
	// 	return err
	// }

	// body := new(bytes.Buffer)
	// if err = tmpl.ExecuteTemplate(body, "body", data); err != nil {
	// 	return err
	// }
	//
	// subject := new(bytes.Buffer)
	// if err = tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
	// 	return err
	// }

	message := gomail.NewMessage()

	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Hello from ZenFollow") // should receive template file that will be constructed

	message.SetBody("text/plain", "test body") // should receive a body

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
