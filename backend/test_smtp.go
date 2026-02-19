//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

func main() {
	// SMTP2Go credentials
	smtpHost := "mail-eu.smtp2go.com"
	smtpPort := "2525"
	smtpUsername := "anytimetravel.app"
	senderEmail := "info@anytimetravel.app"
	senderPassword := os.Getenv("SENDER_PASSWORD")

	if senderPassword == "" {
		fmt.Println("Error: SENDER_PASSWORD environment variable not set")
		fmt.Println("Please set it using: export SENDER_PASSWORD='your_password'")
		return
	}

	recipientEmail := "elias.s.sakr10@gmail.com"

	fmt.Printf("=== SMTP Connection Test ===\n")
	fmt.Printf("SMTP Host: %s\n", smtpHost)
	fmt.Printf("SMTP Port: %s\n", smtpPort)
	fmt.Printf("SMTP Username: %s\n", smtpUsername)
	fmt.Printf("Sender: %s\n", senderEmail)
	fmt.Printf("Recipient: %s\n", recipientEmail)
	fmt.Println("\nAttempting to send test email...")

	// Create authentication
	auth := smtp.PlainAuth("", smtpUsername, senderPassword, smtpHost)

	// Prepare email message
	subject := "Test Email from Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body>
  <h2>SMTP Test Email</h2>
  <p>This is a test email sent at %s</p>
  <p>If you received this, SMTP is working correctly!</p>
  <p>Best regards,<br/>Anytime Travel Team</p>
</body>
</html>
`, time.Now().Format("2006-01-02 15:04:05"))

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		senderEmail, recipientEmail, subject, body,
	)

	// Send email
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, senderEmail, []string{recipientEmail}, []byte(msg))

	if err != nil {
		fmt.Printf("\n❌ Failed to send email: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Email sent successfully!\n")
	fmt.Printf("Check %s for the test email.\n", recipientEmail)
}
