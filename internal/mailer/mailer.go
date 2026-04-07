package mailer

import "embed"

var FS embed.FS

type Mailer interface {
	Send(name, email string, data any) error
}
