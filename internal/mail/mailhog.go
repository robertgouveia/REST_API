package mail

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"text/template"
	"time"

	"github.com/robertgouveia/social/internal/env"
)

type MailHog struct {
	fromEmail string
	smtpHost  string
	smtpPort  string
	password  string
}

func NewMailHog(fromEmail string) *MailHog {
	return &MailHog{
		fromEmail: fromEmail,
		smtpHost:  env.GetString("SMTP_HOST", "localhost"),
		smtpPort:  env.GetString("SMTP_PORT", "1025"),
		password:  env.GetString("SMTP_PASS", ""),
	}
}

func (m *MailHog) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := m.fromEmail
	to := []string{username}

	//template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := append(subject.Bytes(), []byte("\n")...)
	message = append(message, body.Bytes()...)

	auth := smtp.PlainAuth("", from, m.password, m.smtpHost)

	var retryErr error
	for i := 0; i < MaxRetries; i++ {
		retryErr = smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, from, to, message)
		if retryErr != nil {
			//exponential backoff
			time.Sleep(time.Second * time.Duration(i+1)) // more retries wait longer
			continue
		}

		return http.StatusOK, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts, error: %v", MaxRetries, retryErr)
}
