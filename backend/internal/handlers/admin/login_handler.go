package admin

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"travel/backend/core/utils"
	"travel/backend/internal/models/admin"

	"github.com/gofiber/fiber/v2"
)

// Login renders the admin login page
func (h *AdminHandler) Login(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "login", "full-page", "login.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Login",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetLoginFragment serves the login fragment for HTMX partial loads
func (h *AdminHandler) GetLoginFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "login", "fragments", "login-frag.html")))
}

// HandleLogin verifies admin credentials and issues a JWT cookie
func (h *AdminHandler) HandleLogin(c *fiber.Ctx) error {
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	password := c.FormValue("password")

	// Basic validation (do not reveal password policy details)
	var errs []string
	if msg := utils.ValidateEmail("Email", email); msg != "" {
		errs = append(errs, msg)
	}
	if strings.TrimSpace(password) == "" {
		errs = append(errs, "Password is required")
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	log.Printf("[LOGIN] Attempting login for email: %s", email)

	if email == "" || password == "" {
		log.Printf("[LOGIN] Missing credentials")
		return c.Status(fiber.StatusBadRequest).SendString("Email and password are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check whether this IP or email is currently blocked due to repeated failures
	ipKey := "ip:" + c.IP()
	emailKey := "email:" + email
	if h.loginAttemptRepo != nil {
		if blocked, until, err := h.loginAttemptRepo.IsBlocked(ctx, ipKey); err == nil && blocked {
			log.Printf("[LOGIN] Blocked IP %s until %v", c.IP(), until)
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many login attempts. Please try again later.")
		}
		if email != "" {
			if blocked, until, err := h.loginAttemptRepo.IsBlocked(ctx, emailKey); err == nil && blocked {
				log.Printf("[LOGIN] Blocked email %s until %v", email, until)
				return c.Status(fiber.StatusTooManyRequests).SendString("Too many login attempts. Please try again later.")
			}
		}
	}

	adminUser, err := h.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("[LOGIN] User not found for email %s: %v", email, err)
		// Record failed attempt and potentially block
		if h.loginAttemptRepo != nil {
			_ = h.loginAttemptRepo.RecordAttempt(ctx, ipKey)
			if email != "" {
				_ = h.loginAttemptRepo.RecordAttempt(ctx, emailKey)
			}
			// Count recent attempts within 1 minute
			since := time.Now().Add(-1 * time.Minute)
			if cnt, err := h.loginAttemptRepo.CountAttemptsSince(ctx, ipKey, since); err == nil && cnt >= 3 {
				_ = h.loginAttemptRepo.CreateBlock(ctx, ipKey, time.Now().Add(5*time.Minute))
				log.Printf("[LOGIN] Created block for IP %s due to %d attempts", c.IP(), cnt)
			}
			if email != "" {
				if cnt, err := h.loginAttemptRepo.CountAttemptsSince(ctx, emailKey, since); err == nil && cnt >= 3 {
					_ = h.loginAttemptRepo.CreateBlock(ctx, emailKey, time.Now().Add(5*time.Minute))
					log.Printf("[LOGIN] Created block for email %s due to %d attempts", email, cnt)
				}
			}
		}
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	log.Printf("[LOGIN] Found user, checking password...")
	if !utils.CheckPasswordHash(password, adminUser.Password) {
		log.Printf("[LOGIN] Password check failed")
		// Record failed attempt and potentially block
		if h.loginAttemptRepo != nil {
			_ = h.loginAttemptRepo.RecordAttempt(ctx, ipKey)
			if email != "" {
				_ = h.loginAttemptRepo.RecordAttempt(ctx, emailKey)
			}
			since := time.Now().Add(-1 * time.Minute)
			if cnt, err := h.loginAttemptRepo.CountAttemptsSince(ctx, ipKey, since); err == nil && cnt >= 3 {
				_ = h.loginAttemptRepo.CreateBlock(ctx, ipKey, time.Now().Add(5*time.Minute))
				log.Printf("[LOGIN] Created block for IP %s due to %d attempts", c.IP(), cnt)
			}
			if email != "" {
				if cnt, err := h.loginAttemptRepo.CountAttemptsSince(ctx, emailKey, since); err == nil && cnt >= 3 {
					_ = h.loginAttemptRepo.CreateBlock(ctx, emailKey, time.Now().Add(5*time.Minute))
					log.Printf("[LOGIN] Created block for email %s due to %d attempts", email, cnt)
				}
			}
		}
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	log.Printf("[LOGIN] Password verified, generating token...")
	token, err := utils.GenerateJWT(adminUser.Email, adminUser.Role, h.jwtSecret, 24*time.Hour)
	if err != nil {
		log.Printf("[LOGIN] Token generation failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create token")
	}

	// Set auth cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	log.Printf("[LOGIN] Login successful, redirecting to dashboard")
	return c.Redirect("/admin/dashboard")
}

// HandleLogout clears the admin authentication cookie and redirects to login
func (h *AdminHandler) HandleLogout(c *fiber.Ctx) error {
	// Clear the admin token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour), // Set to past time to expire immediately
	})

	// Clear any other auth token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	log.Printf("[LOGOUT] User logged out successfully")

	// Check if it's an HTMX request
	if c.Get("HX-Request") == "true" {
		// Return HX-Redirect header for HTMX
		c.Set("HX-Redirect", "/admin/login")
		return c.SendStatus(fiber.StatusOK)
	}

	// Standard redirect for non-HTMX requests
	return c.Redirect("/admin/login")
}

// Forget renders the admin forgot password page
func (h *AdminHandler) Forget(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "login", "full-page", "forget.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Forgot Password",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetForgetFragment serves the forget password fragment for HTMX partial loads
func (h *AdminHandler) GetForgetFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "login", "fragments", "forget-frag.html")))
}

// HandleForgotPassword handles the forgot password form submission
func (h *AdminHandler) HandleForgotPassword(c *fiber.Ctx) error {
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))

	var errs []string
	if msg := utils.ValidateEmail("Email", email); msg != "" {
		errs = append(errs, msg)
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user exists with this email
	adminUser, err := h.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists for security reasons
		log.Printf("[FORGOT] No user found for email: %s", email)
		return c.JSON(fiber.Map{
			"message": "If an account exists with this email, you will receive a password reset link",
		})
	}

	// Generate 6-digit numeric code
	code := generateNumericCode(6)
	log.Printf("[FORGOT] Generated verification code for %s", email)

	// Create password reset record with code (valid for 15 minutes)
	id := generateSecureToken(16)
	passwordReset := &admin.PasswordReset{
		ID:        id,
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		Used:      false,
	}

	// Store code in database
	if err := h.passwordResetRepo.Create(ctx, passwordReset); err != nil {
		log.Printf("[FORGOT] Failed to store reset code: %v", err)
		return c.JSON(fiber.Map{"message": "If an account exists with this email, you will receive a verification code"})
	}

	// Send email with verification code
	emailService := utils.NewEmailService()
	if err := emailService.SendPasswordResetCodeEmail(adminUser.Email, code); err != nil {
		log.Printf("[FORGOT] Failed to send reset code to %s: %v", adminUser.Email, err)
	}

	// Redirect user to verify page (do not reveal whether email exists)
	return c.JSON(fiber.Map{"message": "If an account exists with this email, you will receive a verification code", "redirect": "/admin/verify?email=" + urlEncode(email)})
}

// generateSecureToken generates a random secure token
func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, err := crand.Read(bytes)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

// generateNumericCode creates a numeric code of given length
func generateNumericCode(length int) string {
	max := 1
	for i := 0; i < length; i++ {
		max *= 10
	}
	n, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Printf("Error generating numeric code: %v", err)
		return "000000"
	}
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n.Int64())
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

// Verify renders the email verification page
func (h *AdminHandler) Verify(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "login", "full-page", "verify.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Verify Email",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetVerifyFragment serves the verification fragment for HTMX partial loads
func (h *AdminHandler) GetVerifyFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "login", "fragments", "verify-frag.html")))
}

// HandleVerify handles the verification code submission
func (h *AdminHandler) HandleVerify(c *fiber.Ctx) error {
	email := c.FormValue("email")
	code := c.FormValue("code")

	var errs []string
	if msg := utils.ValidateEmail("Email", email); msg != "" {
		errs = append(errs, msg)
	}
	if strings.TrimSpace(code) == "" || len(code) != 6 || !isDigits(code) {
		errs = append(errs, "Verification code must be 6 digits")
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the code record
	resetRecord, err := h.passwordResetRepo.FindByCode(ctx, email, code)
	if err != nil {
		log.Printf("[VERIFY] Code verification failed for %s: %v", email, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired verification code"})
	}

	// Generate a token to allow password reset and update the record
	token := generateSecureToken(32)
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := h.passwordResetRepo.UpdateTokenByID(ctx, resetRecord.ID, token, expiresAt); err != nil {
		log.Printf("[VERIFY] Failed to update reset token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process verification"})
	}

	// Respond with redirect to reset page containing token
	return c.JSON(fiber.Map{"message": "Verification successful", "redirect": "/admin/reset?token=" + token})
}

// Reset renders the reset password page
func (h *AdminHandler) Reset(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "login", "full-page", "reset.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Reset Password",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetResetFragment serves the reset password fragment for HTMX partial loads
func (h *AdminHandler) GetResetFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "login", "fragments", "reset-frag.html")))
}

// HandleResetPassword handles the password reset form submission
func (h *AdminHandler) HandleResetPassword(c *fiber.Ctx) error {
	token := c.Query("token", c.FormValue("token"))
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	var errs []string
	if strings.TrimSpace(token) == "" {
		errs = append(errs, "Invalid or missing reset token")
	}
	if msg := utils.ValidateLengthBetween("Password", password, 8, 72); msg != "" {
		errs = append(errs, msg)
	}
	if strings.TrimSpace(confirmPassword) == "" {
		errs = append(errs, "Confirm password is required")
	}
	if password != confirmPassword {
		errs = append(errs, "Passwords do not match")
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find and validate reset token
	resetRecord, err := h.passwordResetRepo.FindByToken(ctx, token)
	if err != nil {
		log.Printf("[RESET] Invalid or expired token: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired reset token",
		})
	}

	// Find the admin user
	adminUser, err := h.adminRepo.FindByEmail(ctx, resetRecord.Email)
	if err != nil {
		log.Printf("[RESET] User not found: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reset password",
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("[RESET] Failed to hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reset password",
		})
	}

	// Update password in database
	adminUser.Password = hashedPassword
	err = h.adminRepo.UpdateByEmail(ctx, adminUser.Email, adminUser)
	if err != nil {
		log.Printf("[RESET] Failed to update password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reset password",
		})
	}

	// Mark token as used
	if err := h.passwordResetRepo.MarkAsUsed(ctx, token); err != nil {
		log.Printf("[RESET] Failed to mark token as used: %v", err)
		// Don't fail the reset even if we can't mark token
	}

	log.Printf("[RESET] Password successfully reset for %s", adminUser.Email)
	return c.JSON(fiber.Map{
		"message":  "Password reset successfully. You can now login with your new password.",
		"redirect": "/admin/login",
	})
}

// isDigits reports whether the entire string is numeric.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
