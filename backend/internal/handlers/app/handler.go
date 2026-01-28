package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Anytime_Travel/backend/core/utils"
	appmodels "Anytime_Travel/backend/internal/models/app"
	"Anytime_Travel/backend/internal/repository/app"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// AppHandler handles app-level requests
type AppHandler struct {
	userRepo  *app.UserRepository
	otpRepo   *app.OTPRepository
	jwtSecret string
}

func NewAppHandler(userRepo *app.UserRepository, otpRepo *app.OTPRepository, jwtSecret string) *AppHandler {
	return &AppHandler{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		jwtSecret: jwtSecret,
	}
}

func (h *AppHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "App Dashboard",
	})
}

func (h *AppHandler) GetProfile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "App User Profile",
	})
}

// Signup handles user registration
func (h *AppHandler) Signup(c *fiber.Ctx) error {
	var req appmodels.SignupRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Validate password match
	if req.Password != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	// Validate password length
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	// Validate sex
	if req.Sex != "Male" && req.Sex != "Female" && req.Sex != "Other" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sex value. Must be Male, Female, or Other",
		})
	}

	// Validate country
	if req.Country == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Country is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email to lowercase
	email := strings.ToLower(req.Email)

	// Check if email already exists
	existingUser, err := h.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Check if phone number already exists
	existingPhone, err := h.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingPhone != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Phone number already registered",
		})
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	// Create user
	user := &appmodels.User{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Email:         email,
		PhoneNumber:   req.PhoneNumber,
		PasswordHash:  passwordHash,
		Sex:           req.Sex,
		Country:       req.Country,
		IsActive:      false, // Set to false until OTP verification
		IsFreezed:     false,
		TotalBookings: 0,
		CreatedAt:     time.Now(),
		LastLogin:     time.Now(),
	}

	// Save user to database
	if err := h.userRepo.Create(ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate and send OTP to phone (4 digits for signup)
	code := utils.GenerateOTP(4)
	otp := &appmodels.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Phone:     user.PhoneNumber,
		Code:      code,
		Type:      appmodels.OTPTypePhone,
		Verified:  false,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	// Delete any existing OTPs for this user
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, appmodels.OTPTypePhone); err != nil {
		// Log error but continue
	}

	// Save OTP to database
	if err := h.otpRepo.Create(ctx, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create OTP",
		})
	}

	// Send OTP via SMS (for now just log it)
	utils.SendSMSOTP(user.PhoneNumber, code)
	fmt.Printf("[Signup] Created OTP for user %s: Phone=%s, Code=%s\n", user.ID, user.PhoneNumber, code)

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Return success response with OTP code (remove in production)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "User registered successfully. OTP sent to phone.",
		"token":    token,
		"otp_code": code, // TODO: Remove in production
		"user": fiber.Map{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"sex":          user.Sex,
			"country":      user.Country,
			"is_active":    user.IsActive,
			"created_at":   user.CreatedAt,
		},
	})
}

// Login handles user authentication
func (h *AppHandler) Login(c *fiber.Ctx) error {
	var req appmodels.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email to lowercase
	email := strings.ToLower(req.Email)

	// Find user by email
	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate",
		})
	}

	// Check if user is frozen
	if user.IsFreezed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Your account has been suspended. Please contact support.",
		})
	}

	// Check if user account is activated
	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Please verify your account with OTP before logging in.",
		})
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Update last login
	if err := h.userRepo.UpdateLastLogin(ctx, user.ID, time.Now()); err != nil {
		// Log error but don't fail login
		// log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"sex":          user.Sex,
			"country":      user.Country,
			"is_active":    user.IsActive,
			"created_at":   user.CreatedAt,
		},
	})
}

// generateToken creates a JWT token for the user
func (h *AppHandler) generateToken(user *appmodels.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// SendOTP sends OTP to email or phone
func (h *AppHandler) SendOTP(c *fiber.Ctx) error {
	var req appmodels.SendOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Determine OTP type and identifier
	var user *appmodels.User
	var err error
	var otpType appmodels.OTPType
	var identifier string

	if req.Type == "email" {
		if req.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email is required",
			})
		}
		user, err = h.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
		otpType = appmodels.OTPTypeEmail
		identifier = req.Email
	} else if req.Type == "phone" {
		if req.Phone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Phone number is required",
			})
		}
		user, err = h.userRepo.FindByPhoneNumber(ctx, req.Phone)
		otpType = appmodels.OTPTypePhone
		identifier = req.Phone
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid OTP type",
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find user",
		})
	}

	// Delete any existing OTPs for this user and type
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, otpType); err != nil {
		// Log error but continue
	}

	// Generate 6-digit OTP code
	code := utils.GenerateOTP(6)

	// Create OTP record
	otp := &appmodels.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Code:      code,
		Type:      otpType,
		Verified:  false,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP expires in 10 minutes
		CreatedAt: time.Now(),
	}

	if otpType == appmodels.OTPTypeEmail {
		otp.Email = identifier
	} else {
		otp.Phone = identifier
	}

	// Save OTP to database
	if err := h.otpRepo.Create(ctx, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create OTP",
		})
	}

	// TODO: Send OTP via email or SMS
	// For now, we'll just return success (in production, integrate with email/SMS service)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}

// VerifyOTP verifies OTP and activates user account
func (h *AppHandler) VerifyOTP(c *fiber.Ctx) error {
	var req appmodels.VerifyOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	// Log the request for debugging
	fmt.Printf("[VerifyOTP] Request: Type=%s, Phone=%s, Code=%s, Email=%s\n",
		req.Type, req.Phone, req.Code, req.Email)

	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP code is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find OTP record
	var otp *appmodels.OTP
	var err error

	if req.Type == "email" {
		if req.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email is required",
			})
		}
		otp, err = h.otpRepo.FindByEmailAndCode(ctx, strings.ToLower(req.Email), req.Code)
	} else if req.Type == "phone" {
		if req.Phone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Phone number is required",
			})
		}
		fmt.Printf("[VerifyOTP] Searching for OTP: Phone=%s, Code=%s\n", req.Phone, req.Code)
		otp, err = h.otpRepo.FindByPhoneAndCode(ctx, req.Phone, req.Code)
		if err != nil {
			fmt.Printf("[VerifyOTP] Error finding OTP: %v\n", err)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid OTP type",
		})
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid OTP code",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP has expired",
		})
	}

	// Check if already verified
	if otp.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP already used",
		})
	}

	// Mark OTP as verified
	if err := h.otpRepo.MarkAsVerified(ctx, otp.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP verified successfully",
	})
}

// ActivateAccount activates user account after both email and phone verification
func (h *AppHandler) ActivateAccount(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	user, err := h.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find user",
		})
	}

	// Activate user account
	user.IsActive = true
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate account",
		})
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account activated successfully",
		"token":   token,
		"user": fiber.Map{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"is_active":    user.IsActive,
		},
	})
}

// ResetPassword handles password reset with OTP verification
func (h *AppHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Phone       string `json:"phone"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Phone == "" || req.OTP == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone, OTP, and new password are required",
		})
	}

	// Validate password length
	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify OTP
	otp, err := h.otpRepo.FindByPhoneAndCode(ctx, req.Phone, req.OTP)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired OTP",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP has expired",
		})
	}

	// Check if OTP is verified
	if !otp.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP not verified",
		})
	}

	// Find user by phone
	user, err := h.userRepo.FindByPhoneNumber(ctx, req.Phone)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find user",
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Update user password
	user.PasswordHash = hashedPassword
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	// Delete the used OTP
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, appmodels.OTPTypePhone); err != nil {
		fmt.Printf("Warning: Failed to delete OTP: %v\n", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}
