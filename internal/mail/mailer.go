package mail

import "embed"

const (
	FromName            = "GopherSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(string, string, string, any, bool) (int, error)
}
