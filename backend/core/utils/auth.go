package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AdminClaims represents JWT claims for admin authentication.
type AdminClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

// GenerateJWT creates a signed JWT with the given subject and role.
func GenerateJWT(subject, role, secret string, ttl time.Duration) (string, error) {
	claims := AdminClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT validates and parses a JWT string into AdminClaims.
func ParseJWT(tokenStr, secret string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// GenerateOTP generates a random n-digit OTP code
func GenerateOTP(length int) string {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		otp[i] = digits[num.Int64()]
	}

	return string(otp)
}

// FormatOTPCode formats OTP for display (e.g., "1234" -> "1 2 3 4")
func FormatOTPCode(code string) string {
	result := ""
	for i, c := range code {
		if i > 0 {
			result += " "
		}
		result += string(c)
	}
	return result
}

// SendEmailOTP sends OTP via email using SMTP.
// Reads SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM, SMTP_FROM_NAME
// from env. If SMTP_HOST is unset, falls back to logging the OTP to stdout (useful for local dev).
func SendEmailOTP(email, code string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		fmt.Printf("[DEV] SMTP not configured; OTP %s would be emailed to %s\n", code, email)
		return nil
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = username
	}
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "Anytime Trip"
	}

	subject := "Your verification code"
	htmlBody := fmt.Sprintf(`<!doctype html><html><body style="font-family:Arial,sans-serif;background:#f5f5f5;padding:24px;">
<div style="max-width:480px;margin:0 auto;background:#ffffff;border-radius:8px;padding:32px;">
<h2 style="color:#1e5a8e;margin:0 0 16px;">Verify your account</h2>
<p style="color:#333;font-size:14px;line-height:1.5;">Use the code below to finish signing in. It expires in 10 minutes.</p>
<div style="font-size:32px;font-weight:bold;letter-spacing:8px;color:#1e5a8e;text-align:center;padding:24px 0;">%s</div>
<p style="color:#777;font-size:12px;">If you didn't request this code, you can safely ignore this email.</p>
</div></body></html>`, code)

	msg := "From: " + fmt.Sprintf("%s <%s>", fromName, from) + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + htmlBody

	addr := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	// Port 465 requires implicit TLS; 587 uses STARTTLS via smtp.SendMail.
	if port == "465" {
		tlsConfig := &tls.Config{ServerName: host}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("smtp tls dial: %w", err)
		}
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("smtp client: %w", err)
		}
		defer client.Quit()
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err := client.Mail(from); err != nil {
			return fmt.Errorf("smtp mail: %w", err)
		}
		if err := client.Rcpt(email); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
		wc, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data: %w", err)
		}
		if _, err := wc.Write([]byte(msg)); err != nil {
			return fmt.Errorf("smtp write: %w", err)
		}
		if err := wc.Close(); err != nil {
			return fmt.Errorf("smtp close: %w", err)
		}
		return nil
	}

	if err := smtp.SendMail(addr, auth, from, []string{email}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

// SendSMSOTP is a stub: no real SMS provider is wired up.
// We log the code so it can be read from server output during a demo if needed.
func SendSMSOTP(phone, code string) error {
	fmt.Printf("[DEV] SMS provider not configured; OTP %s would be sent to %s\n", code, phone)
	return nil
}
