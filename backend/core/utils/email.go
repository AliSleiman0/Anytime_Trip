package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost       string
	smtpPort       string
	senderEmail    string
	senderPassword string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:       getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:       getEnvOrDefault("SMTP_PORT", "587"),
		senderEmail:    getEnvOrDefault("SENDER_EMAIL", "noreply@travel.app"),
		senderPassword: getEnvOrDefault("SENDER_PASSWORD", ""),
	}
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(recipientEmail, resetLink string) error {
	// If no SMTP credentials configured, log and skip (useful for development)
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Password reset link for %s: %s\n", recipientEmail, resetLink)
		return nil
	}

	subject := "Password Reset Request - Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; }
    .header { background-color: #336891; color: white; padding: 20px; text-align: center; }
    .content { padding: 20px; border: 1px solid #ddd; }
    .button { background-color: #336891; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; }
    .footer { text-align: center; color: #666; margin-top: 20px; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Password Reset Request</h1>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>We received a request to reset your password. Click the button below to create a new password:</p>
      <p style="text-align: center; margin: 30px 0;">
        <a href="%s" class="button">Reset Password</a>
      </p>
      <p>This link will expire in 1 hour.</p>
      <p>If you didn't request a password reset, you can safely ignore this email.</p>
      <p>Best regards,<br>Travel Team</p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Travel. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
`, resetLink)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendPasswordResetCodeEmail sends a verification code for password reset
func (es *EmailService) SendPasswordResetCodeEmail(recipientEmail, code string) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Password reset code for %s: %s\n", recipientEmail, code)
		return nil
	}

	subject := "Your password reset code - Travel"
	body := fmt.Sprintf(`
<html>
<body>
  <p>Hello,</p>
  <p>Use the verification code below to reset your password. The code expires in 15 minutes.</p>
  <h2 style="letter-spacing:6px;">%s</h2>
  <p>If you didn't request this, ignore this email.</p>
  <p>Best regards,<br/>Travel Team</p>
</body>
</html>
`, code)

	return es.sendEmail(recipientEmail, subject, body)
}

// sendEmail sends an email using SMTP
func (es *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.senderEmail, es.senderPassword, es.smtpHost)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		es.senderEmail, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
	return smtp.SendMail(addr, auth, es.senderEmail, []string{to}, []byte(msg))
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
