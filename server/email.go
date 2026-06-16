package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/config"
)

// generateVerificationToken creates a cryptographically secure random token
// and sets its expiry to 24 hours from now.
func generateVerificationToken() (string, time.Time, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate verification token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour)
	return token, expiresAt, nil
}

// SendVerificationEmail sends a verification email to the user.
// If SMTP is not configured (host is empty), it logs the verification URL to console.
func SendVerificationEmail(to, token string) error {
	smtpCfg := config.AppConfigInstance.SMTP
	verifyURL := fmt.Sprintf("http://localhost:3000/auth/verify-email?token=%s", token)

	subject := "Подтверждение регистрации — Vera Art"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
		.container { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #6366f1, #8b5cf6); padding: 30px; text-align: center; }
		.header h1 { color: #ffffff; margin: 0; font-size: 24px; }
		.content { padding: 30px; }
		.content p { color: #374151; line-height: 1.6; font-size: 16px; }
		.btn { display: inline-block; padding: 14px 32px; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #ffffff !important; text-decoration: none; border-radius: 8px; font-size: 16px; font-weight: 600; margin: 20px 0; }
		.footer { padding: 20px 30px; background: #f9fafb; text-align: center; font-size: 13px; color: #9ca3af; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Vera Art</h1>
		</div>
		<div class="content">
			<p>Здравствуйте!</p>
			<p>Спасибо за регистрацию на платформе <strong>Vera Art</strong>. Для завершения регистрации, пожалуйста, подтвердите свой email, нажав на кнопку ниже:</p>
			<p style="text-align: center;">
				<a href="%s" class="btn">Подтвердить email</a>
			</p>
			<p>Или скопируйте и вставьте эту ссылку в браузер:</p>
			<p style="word-break: break-all; color: #6366f1; font-size: 14px;">%s</p>
			<p>Ссылка действительна в течение <strong>24 часов</strong>.</p>
			<p>Если вы не регистрировались на Vera Art, просто проигнорируйте это письмо.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Vera Art. Все права защищены.</p>
		</div>
	</div>
</body>
</html>`, verifyURL, verifyURL)

	// If SMTP is not configured, log to console (dev mode)
	if smtpCfg.Host == "" {
		log.Printf("[EMAIL] Verification email for %s", to)
		log.Printf("[EMAIL] Subject: %s", subject)
		log.Printf("[EMAIL] Verification URL: %s", verifyURL)
		log.Printf("[EMAIL] SMTP not configured — email logged to console only")
		return nil
	}

	// Build email message
	msg := buildMIMEMessage(smtpCfg.From, to, subject, body)

	// SMTP auth
	auth := smtp.PlainAuth("", smtpCfg.Username, smtpCfg.Password, smtpCfg.Host)
	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)

	if err := smtp.SendMail(addr, auth, smtpCfg.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("[EMAIL] Verification email sent successfully to %s", to)
	return nil
}

// buildMIMEMessage constructs a proper MIME email message with HTML content.
func buildMIMEMessage(from, to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""
	headers["Content-Transfer-Encoding"] = "base64"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}
