package mail

import "embed"

const (
	FromName            = "GopherSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "template"
var FS embed.FS

type Client interface {
	Send(string, string, string, any, bool) error
}
