package mailer

import "embed"

//go:embed templates
var FS embed.FS

type Mailer interface {
	Send(name, email string, data any) error
}
