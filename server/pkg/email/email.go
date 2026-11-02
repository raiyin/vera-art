package email

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"strings"
)

// Sender handles email sending operations.
type Sender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewSender creates a new email sender.
func NewSender(host string, port int, username, password, from string) *Sender {
	return &Sender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendVerificationEmail sends an email verification message.
func (s *Sender) SendVerificationEmail(to, token string) error {
	subject := "Подтверждение email"
	body := fmt.Sprintf(`
		<h2>Подтверждение email</h2>
		<p>Для подтверждения email перейдите по ссылке:</p>
		<p><a href="%s">Подтвердить email</a></p>
		<p>Если вы не регистрировались на сайте, проигнорируйте это письмо.</p>
	`, token)

	return s.send(to, subject, body)
}

// send sends an email via SMTP.
func (s *Sender) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	mime := buildMIMEMessage(s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(mime))
}

// buildMIMEMessage builds a MIME-formatted email message.
func buildMIMEMessage(from, to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}

// GenerateVerificationToken generates a random verification token.
func GenerateVerificationToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Ensure Sender implements the EmailSender interface.
var _ EmailSender = (*Sender)(nil)

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendVerificationEmail(to, token string) error
}
