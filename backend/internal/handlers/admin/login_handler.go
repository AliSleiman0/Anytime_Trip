package admin

import (
	"context"
	"html/template"
	"log"
	"path/filepath"
	"strings"
	"time"

	"Anytime_Travel/backend/core/utils"

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

	adminUser, err := h.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("[LOGIN] User not found for email %s: %v", email, err)
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	log.Printf("[LOGIN] Found user, checking password...")
	if !utils.CheckPasswordHash(password, adminUser.Password) {
		log.Printf("[LOGIN] Password check failed")
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
	email := c.FormValue("email")

	var errs []string
	if msg := utils.ValidateEmail("Email", email); msg != "" {
		errs = append(errs, msg)
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	// TODO: Implement password reset logic:
	// 1. Check if user exists with this email
	// 2. Generate reset token
	// 3. Store token in database with expiration
	// 4. Send email with reset link

	// For now, return success message
	return c.JSON(fiber.Map{
		"message": "If an account exists with this email, you will receive a password reset link",
	})
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

	// TODO: Implement verification logic:
	// 1. Check if verification code matches the one sent to email
	// 2. Check if code has not expired
	// 3. If valid, mark email as verified
	// 4. Redirect to password reset page

	// For now, return success message
	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
	})
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
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	var errs []string
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

	// TODO: Implement password reset logic:
	// 1. Get the user email from session/context
	// 2. Hash the new password
	// 3. Update the password in database
	// 4. Clear any reset tokens for this user
	// 5. Redirect to login page

	// For now, return success message
	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
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
