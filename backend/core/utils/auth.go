package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

// SendEmailOTP sends OTP via email (placeholder for actual email service)
func SendEmailOTP(email, code string) error {
	// TODO: Implement actual email sending logic
	// For now, just log it
	fmt.Printf("Sending OTP %s to email %s\n", code, email)
	return nil
}

// SendSMSOTP sends OTP via SMS (placeholder for actual SMS service)
func SendSMSOTP(phone, code string) error {
	// TODO: Implement actual SMS sending logic
	// For now, just log it
	fmt.Printf("Sending OTP %s to phone %s\n", code, phone)
	return nil
}
