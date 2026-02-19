package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"Anytime_Travel/backend/core/utils"
	adminmodels "Anytime_Travel/backend/internal/models/admin"
	appmodels "Anytime_Travel/backend/internal/models/app"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	"Anytime_Travel/backend/internal/repository/app"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AppHandler handles app-level requests
type AppHandler struct {
	userRepo            *app.UserRepository
	otpRepo             *app.OTPRepository
	paymentMethodRepo   *app.PaymentMethodRepository
	travelRepo          *adminrepo.TravelRepository
	bannerRepo          *adminrepo.BannerRepository
	searchBannerRepo    *adminrepo.SearchBannerRepository
	popularRepo         *adminrepo.PopularRepository
	carRepo             *adminrepo.CarRepository
	hotelBookingRepo    *app.HotelBookingRepository
	flightBookingRepo   *app.FlightBookingRepository
	carBookingRepo      *app.CarBookingRepository
	transferBookingRepo *app.TransferBookingRepository
	supportTicketRepo   *app.SupportTicketRepository
	chatbotRatingRepo   *app.ChatbotRatingRepository
	notificationHelper  *utils.NotificationHelper
	emailService        *utils.EmailService
	jwtSecret           string
}

func NewAppHandler(
	userRepo *app.UserRepository,
	otpRepo *app.OTPRepository,
	paymentMethodRepo *app.PaymentMethodRepository,
	travelRepo *adminrepo.TravelRepository,
	bannerRepo *adminrepo.BannerRepository,
	searchBannerRepo *adminrepo.SearchBannerRepository,
	popularRepo *adminrepo.PopularRepository,
	carRepo *adminrepo.CarRepository,
	hotelBookingRepo *app.HotelBookingRepository,
	flightBookingRepo *app.FlightBookingRepository,
	carBookingRepo *app.CarBookingRepository,
	transferBookingRepo *app.TransferBookingRepository,
	supportTicketRepo *app.SupportTicketRepository,
	chatbotRatingRepo *app.ChatbotRatingRepository,
	notificationHelper *utils.NotificationHelper,
	emailService *utils.EmailService,
	jwtSecret string,
) *AppHandler {
	return &AppHandler{
		userRepo:            userRepo,
		otpRepo:             otpRepo,
		paymentMethodRepo:   paymentMethodRepo,
		travelRepo:          travelRepo,
		bannerRepo:          bannerRepo,
		searchBannerRepo:    searchBannerRepo,
		popularRepo:         popularRepo,
		carRepo:             carRepo,
		hotelBookingRepo:    hotelBookingRepo,
		flightBookingRepo:   flightBookingRepo,
		carBookingRepo:      carBookingRepo,
		transferBookingRepo: transferBookingRepo,
		supportTicketRepo:   supportTicketRepo,
		chatbotRatingRepo:   chatbotRatingRepo,
		notificationHelper:  notificationHelper,
		emailService:        emailService,
		jwtSecret:           jwtSecret,
	}
}

func (h *AppHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "App Dashboard",
	})
}

func (h *AppHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from JWT token (assuming middleware sets it)
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by ID
	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile fetched successfully",
		"user": fiber.Map{
			"id":             user.ID,
			"name":           user.Name,
			"email":          user.Email,
			"phone_number":   user.PhoneNumber,
			"profile_image":  user.ProfileImage,
			"sex":            user.Sex,
			"country":        user.Country,
			"total_bookings": user.TotalBookings,
			"is_active":      user.IsActive,
			"created_at":     user.CreatedAt,
			"last_login":     user.LastLogin,
		},
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

	// Generate and send OTP to email (6 digits for signup)
	code := utils.GenerateOTP(6)
	otp := &appmodels.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Email:     user.Email,
		Code:      code,
		Type:      appmodels.OTPTypeEmail,
		Verified:  false,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	// Delete any existing OTPs for this user
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, appmodels.OTPTypeEmail); err != nil {
		// Log error but continue
	}

	// Save OTP to database
	if err := h.otpRepo.Create(ctx, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create OTP",
		})
	}

	// Send OTP via email
	if err := h.emailService.SendOTPEmail(user.Email, code, "verify your account"); err != nil {
		fmt.Printf("[Signup] Failed to send OTP email: %v\n", err)
		// Continue even if email fails - user can request resend
	}

	fmt.Printf("[Signup] Created OTP for user %s: Email=%s, Code=%s\n", user.ID, user.Email, code)

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Return success response with OTP code (remove in production)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "User registered successfully. OTP sent to email.",
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
		"sub":     user.ID, // Standard JWT subject field
		"user_id": user.ID, // Keep for backward compatibility
		"email":   user.Email,
		"name":    user.Name,
		"role":    "user",                                    // Add role for middleware
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
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
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

	var user *appmodels.User
	var otp *appmodels.OTP
	var err error
	var otpType appmodels.OTPType

	// Prioritize email-based reset over phone-based
	if req.Email != "" && req.Code != "" {
		// Email-based password reset
		email := strings.ToLower(req.Email)

		// Verify OTP
		otp, err = h.otpRepo.FindByEmailAndCode(ctx, email, req.Code)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid or expired verification code",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify code",
			})
		}

		// Find user by email
		user, err = h.userRepo.FindByEmail(ctx, email)
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
		otpType = appmodels.OTPTypeEmail
	} else if req.Phone != "" && req.OTP != "" {
		// Phone-based password reset (legacy support)
		// Verify OTP
		otp, err = h.otpRepo.FindByPhoneAndCode(ctx, req.Phone, req.OTP)
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

		// Find user by phone
		user, err = h.userRepo.FindByPhoneNumber(ctx, req.Phone)
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
		otpType = appmodels.OTPTypePhone
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and code, or phone and OTP are required",
		})
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verification code has expired",
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
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, otpType); err != nil {
		fmt.Printf("Warning: Failed to delete OTP: %v\n", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// ForgotPassword handles the forgot password request by sending OTP via email
func (h *AppHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate email
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email
	email := strings.ToLower(req.Email)

	// Check if user exists
	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No account found with this email",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find user",
		})
	}

	// Delete any existing OTPs for this user's email
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, appmodels.OTPTypeEmail); err != nil {
		fmt.Printf("Warning: Failed to delete existing OTPs: %v\n", err)
	}

	// Generate 6-digit OTP code
	code := utils.GenerateOTP(6)

	// Create OTP record
	otp := &appmodels.OTP{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Email:     email,
		Code:      code,
		Type:      appmodels.OTPTypeEmail,
		Verified:  false,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP expires in 10 minutes
		CreatedAt: time.Now(),
	}

	// Save OTP to database
	if err := h.otpRepo.Create(ctx, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create OTP",
		})
	}

	// Send OTP via email
	if err := h.emailService.SendPasswordResetCodeEmail(email, code); err != nil {
		fmt.Printf("Error sending password reset email: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send password reset email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset code has been sent to your email",
	})
}

// VerifyResetOTP verifies the OTP for password reset
func (h *AppHandler) VerifyResetOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Email == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and code are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email
	email := strings.ToLower(req.Email)

	// Find OTP by email and code
	otp, err := h.otpRepo.FindByEmailAndCode(ctx, email, req.Code)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired OTP code",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP code has expired",
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

// ResetPasswordWithEmail handles password reset using email-based OTP
func (h *AppHandler) ResetPasswordWithEmail(c *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, code, and new password are required",
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

	// Normalize email
	email := strings.ToLower(req.Email)

	// Verify OTP
	otp, err := h.otpRepo.FindByEmailAndCode(ctx, email, req.Code)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired OTP code",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP code has expired",
		})
	}

	// Check if OTP is verified
	if !otp.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "OTP code not verified. Please verify the code first.",
		})
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(ctx, email)
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
	if err := h.otpRepo.DeleteByUserAndType(ctx, user.ID, appmodels.OTPTypeEmail); err != nil {
		fmt.Printf("Warning: Failed to delete OTP: %v\n", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully. You can now login with your new password.",
	})
}

// GoogleSignIn handles Google OAuth sign-in
func (h *AppHandler) GoogleSignIn(c *fiber.Ctx) error {
	fmt.Println("========== GOOGLE SIGN-IN REQUEST RECEIVED ==========")

	var req struct {
		UID      string `json:"uid"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		PhotoURL string `json:"photo_url"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("[ERROR] Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	fmt.Printf("[INFO] Request data: UID=%s, Email=%s, Name=%s\n", req.UID, req.Email, req.Name)

	// Validate required fields
	if req.UID == "" || req.Email == "" {
		fmt.Println("[ERROR] Missing required fields: UID or Email")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UID and email are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email to lowercase
	email := strings.ToLower(req.Email)
	fmt.Printf("[INFO] Normalized email: %s\n", email)

	// Check if user exists by email
	fmt.Println("[INFO] Checking if user exists in database...")
	existingUser, err := h.userRepo.FindByEmail(ctx, email)
	isNewUser := false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// New user - create account
			fmt.Println("[INFO] User not found - creating new user")
			isNewUser = true
			user := &appmodels.User{
				ID:            uuid.New().String(),
				Name:          req.Name,
				Email:         email,
				PhoneNumber:   "", // To be filled in profile completion
				PasswordHash:  "", // No password for Google sign-in
				Sex:           "", // To be filled in profile completion
				Country:       "", // To be filled in profile completion
				IsActive:      true,
				IsFreezed:     false,
				TotalBookings: 0,
				CreatedAt:     time.Now(),
				LastLogin:     time.Now(),
			}

			fmt.Printf("[INFO] Creating new user with ID: %s\n", user.ID)

			// Save user to database
			if err := h.userRepo.Create(ctx, user); err != nil {
				fmt.Printf("[ERROR] Failed to create user in database: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user",
				})
			}

			fmt.Println("[SUCCESS] New user created successfully")
			existingUser = user
		} else {
			fmt.Printf("[ERROR] Database error while finding user: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	} else {
		// Existing user - update last login
		fmt.Printf("[INFO] Existing user found: %s (ID: %s)\n", existingUser.Email, existingUser.ID)
		existingUser.LastLogin = time.Now()
		if err := h.userRepo.Update(ctx, existingUser.ID, existingUser); err != nil {
			fmt.Printf("[WARNING] Failed to update last login: %v\n", err)
		} else {
			fmt.Println("[INFO] Last login updated successfully")
		}
	}

	// Generate JWT token
	fmt.Println("[INFO] Generating JWT token...")
	tokenString, err := utils.GenerateJWT(existingUser.ID, "user", h.jwtSecret, 24*7*time.Hour)
	if err != nil {
		fmt.Printf("[ERROR] Failed to generate JWT token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	fmt.Printf("[SUCCESS] JWT token generated successfully\n")
	fmt.Printf("[INFO] Returning response - isNewUser: %v\n", isNewUser)
	fmt.Println("========== GOOGLE SIGN-IN COMPLETED ==========")

	// Return user data and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Sign-in successful",
		"token":     tokenString,
		"isNewUser": isNewUser,
		"user": fiber.Map{
			"id":             existingUser.ID,
			"name":           existingUser.Name,
			"email":          existingUser.Email,
			"phone_number":   existingUser.PhoneNumber,
			"profile_image":  existingUser.ProfileImage,
			"sex":            existingUser.Sex,
			"country":        existingUser.Country,
			"total_bookings": existingUser.TotalBookings,
			"is_active":      existingUser.IsActive,
			"created_at":     existingUser.CreatedAt,
			"last_login":     existingUser.LastLogin,
		},
	})
}

// AppleSignIn handles Apple sign-in/signup for app users
func (h *AppHandler) AppleSignIn(c *fiber.Ctx) error {
	fmt.Println("========== APPLE SIGN-IN REQUEST RECEIVED ==========")

	var req struct {
		UID         string `json:"uid"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		PhotoURL    string `json:"photo_url"`
		AppleUserID string `json:"apple_user_id"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("[ERROR] Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	fmt.Printf("[INFO] Request data: UID=%s, Email=%s, Name=%s, AppleUserID=%s\n", req.UID, req.Email, req.Name, req.AppleUserID)

	// Validate required fields
	if req.UID == "" {
		fmt.Println("[ERROR] Missing required field: UID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Normalize email to lowercase if provided
	email := ""
	if req.Email != "" {
		email = strings.ToLower(req.Email)
		fmt.Printf("[INFO] Normalized email: %s\n", email)
	}

	// Check if user exists by email (if email is available)
	fmt.Println("[INFO] Checking if user exists in database...")
	var existingUser *appmodels.User
	var err error
	isNewUser := false

	if email != "" {
		existingUser, err = h.userRepo.FindByEmail(ctx, email)
	} else {
		// Apple sign-in may not provide email, so we mark it as new user
		err = mongo.ErrNoDocuments
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// New user - create account
			fmt.Println("[INFO] User not found - creating new user")
			isNewUser = true

			// Use Apple User ID as email if email is not provided
			userEmail := email
			if userEmail == "" {
				userEmail = req.AppleUserID + "@privaterelay.appleid.com"
			}

			user := &appmodels.User{
				ID:            uuid.New().String(),
				Name:          req.Name,
				Email:         userEmail,
				PhoneNumber:   "", // To be filled in profile completion
				PasswordHash:  "", // No password for Apple sign-in
				Sex:           "", // To be filled in profile completion
				Country:       "", // To be filled in profile completion
				IsActive:      true,
				IsFreezed:     false,
				TotalBookings: 0,
				CreatedAt:     time.Now(),
				LastLogin:     time.Now(),
			}

			fmt.Printf("[INFO] Creating new user with ID: %s\n", user.ID)

			// Save user to database
			if err := h.userRepo.Create(ctx, user); err != nil {
				fmt.Printf("[ERROR] Failed to create user in database: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user",
				})
			}

			fmt.Println("[SUCCESS] New user created successfully")
			existingUser = user
		} else {
			fmt.Printf("[ERROR] Database error while finding user: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	} else {
		// Existing user - update last login
		fmt.Printf("[INFO] Existing user found: %s (ID: %s)\n", existingUser.Email, existingUser.ID)
		existingUser.LastLogin = time.Now()
		if err := h.userRepo.Update(ctx, existingUser.ID, existingUser); err != nil {
			fmt.Printf("[WARNING] Failed to update last login: %v\n", err)
		} else {
			fmt.Println("[INFO] Last login updated successfully")
		}
	}

	// Generate JWT token
	fmt.Println("[INFO] Generating JWT token...")
	tokenString, err := utils.GenerateJWT(existingUser.ID, "user", h.jwtSecret, 24*7*time.Hour)
	if err != nil {
		fmt.Printf("[ERROR] Failed to generate JWT token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	fmt.Printf("[SUCCESS] JWT token generated successfully\n")
	fmt.Printf("[INFO] Returning response - isNewUser: %v\n", isNewUser)
	fmt.Println("========== APPLE SIGN-IN COMPLETED ==========")

	// Return user data and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Sign-in successful",
		"token":     tokenString,
		"isNewUser": isNewUser,
		"user": fiber.Map{
			"id":             existingUser.ID,
			"name":           existingUser.Name,
			"email":          existingUser.Email,
			"phone_number":   existingUser.PhoneNumber,
			"profile_image":  existingUser.ProfileImage,
			"sex":            existingUser.Sex,
			"country":        existingUser.Country,
			"total_bookings": existingUser.TotalBookings,
			"is_active":      existingUser.IsActive,
			"created_at":     existingUser.CreatedAt,
			"last_login":     existingUser.LastLogin,
		},
	})
}

// UpdateProfile updates user profile information
func (h *AppHandler) UpdateProfile(c *fiber.Ctx) error {
	fmt.Println("========== UPDATE PROFILE REQUEST ==========")

	var req struct {
		PhoneNumber string `json:"phone_number"`
		Sex         string `json:"sex"`
		Country     string `json:"country"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("[ERROR] Failed to parse request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	fmt.Printf("[INFO] Update data - Phone: %s, Sex: %s, Country: %s\n", req.PhoneNumber, req.Sex, req.Country)

	// Get user ID from JWT token (assuming middleware sets it)
	userID := c.Locals("user_id")
	fmt.Printf("[INFO] User ID from token: %v\n", userID)

	if userID == nil {
		fmt.Println("[ERROR] No user_id in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	fmt.Printf("[INFO] Looking up user with ID: %s\n", userID.(string))
	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to find user: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	fmt.Printf("[INFO] Found user: %s (%s)\n", user.Email, user.ID)

	// Update fields if provided
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.Sex != "" {
		user.Sex = req.Sex
	}
	if req.Country != "" {
		user.Country = req.Country
	}

	// Save updated user
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		fmt.Printf("[ERROR] Failed to update user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	fmt.Println("[SUCCESS] Profile updated successfully")
	fmt.Println("========== UPDATE PROFILE COMPLETED ==========")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":             user.ID,
			"name":           user.Name,
			"email":          user.Email,
			"phone_number":   user.PhoneNumber,
			"profile_image":  user.ProfileImage,
			"sex":            user.Sex,
			"country":        user.Country,
			"total_bookings": user.TotalBookings,
			"is_active":      user.IsActive,
			"created_at":     user.CreatedAt,
			"last_login":     user.LastLogin,
		},
	})
}

// UploadProfileImage handles profile image upload
func (h *AppHandler) UploadProfileImage(c *fiber.Ctx) error {
	fmt.Println("========== UPLOAD PROFILE IMAGE REQUEST ==========")

	// Get user ID from JWT token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Printf("[ERROR] Failed to get file: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size exceeds 5MB limit",
		})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Only images are allowed",
		})
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./static/uploads/profiles"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("[ERROR] Failed to create upload directory: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", userID.(string), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file
	if err := c.SaveFile(file, filePath); err != nil {
		fmt.Printf("[ERROR] Failed to save file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image",
		})
	}

	// Generate public URL
	imageURL := fmt.Sprintf("/uploads/profiles/%s", filename)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to find user: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Delete old profile image if exists
	if user.ProfileImage != "" {
		oldPath := strings.TrimPrefix(user.ProfileImage, "/")
		oldPath = filepath.Join("./static", oldPath)
		if err := os.Remove(oldPath); err != nil {
			fmt.Printf("[WARN] Failed to delete old image: %v\n", err)
		}
	}

	// Update user's profile image
	user.ProfileImage = imageURL
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		fmt.Printf("[ERROR] Failed to update user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile image",
		})
	}

	fmt.Printf("[SUCCESS] Profile image uploaded: %s\n", imageURL)
	fmt.Println("========== UPLOAD PROFILE IMAGE COMPLETED ==========")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Profile image uploaded successfully",
		"profile_image": imageURL,
	})
}

// GetNotificationPreferences retrieves the user's notification preferences
func (h *AppHandler) GetNotificationPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Notification preferences retrieved successfully",
		"preferences": user.NotificationPreferences,
	})
}

// SaveNotificationPreferences saves or updates the user's notification preferences
func (h *AppHandler) SaveNotificationPreferences(c *fiber.Ctx) error {
	fmt.Println("========== SAVE NOTIFICATION PREFERENCES REQUEST ==========")

	userID := c.Locals("user_id")
	if userID == nil {
		fmt.Println("[ERROR] No user_id in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var prefs appmodels.NotificationPreferences
	if err := c.BodyParser(&prefs); err != nil {
		fmt.Printf("[ERROR] Failed to parse request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to find user: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update notification preferences
	user.NotificationPreferences = prefs

	// Save updated user
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		fmt.Printf("[ERROR] Failed to update user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save notification preferences",
		})
	}

	fmt.Println("[SUCCESS] Notification preferences saved successfully")
	fmt.Println("========== SAVE NOTIFICATION PREFERENCES COMPLETED ==========")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Notification preferences saved successfully",
		"preferences": user.NotificationPreferences,
	})
}

// AddPaymentMethod adds a new payment method for the user
func (h *AppHandler) AddPaymentMethod(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req appmodels.AddPaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Extract only last 4 digits - NEVER store full card number or CVV
	var cardLast4 string
	if len(req.CardNumber) >= 4 {
		cardLast4 = req.CardNumber[len(req.CardNumber)-4:]
	}

	pm := &appmodels.PaymentMethod{
		UserID:          userID.(string),
		CardHolderName:  req.CardHolderName,
		CardNumberLast4: cardLast4,
		CardBrand:       req.CardBrand,
		ExpiryMonth:     req.ExpiryMonth,
		ExpiryYear:      req.ExpiryYear,
		BillingAddress:  req.BillingAddress,
		City:            req.City,
		State:           req.State,
		PostalCode:      req.PostalCode,
		Country:         req.Country,
		PhoneNumber:     req.PhoneNumber,
		IsDefault:       req.IsDefault,
		IsActive:        true,
	}

	result, err := h.paymentMethodRepo.AddPaymentMethod(ctx, pm)
	if err != nil {
		fmt.Printf("[ERROR] Failed to add payment method: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add payment method",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payment method added successfully",
		"data":    result,
	})
}

// GetPaymentMethods retrieves all payment methods for the user
func (h *AppHandler) GetPaymentMethods(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	methods, err := h.paymentMethodRepo.GetPaymentMethodsByUserID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to retrieve payment methods: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve payment methods",
		})
	}

	if methods == nil {
		methods = []appmodels.PaymentMethod{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment methods retrieved successfully",
		"data":    methods,
	})
}

// UpdatePaymentMethod updates a payment method
func (h *AppHandler) UpdatePaymentMethod(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment method ID is required",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment method ID",
		})
	}

	var req appmodels.UpdatePaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Get existing payment method
	existing, err := h.paymentMethodRepo.GetPaymentMethodByID(ctx, objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment method not found",
		})
	}

	// Verify ownership
	if existing.UserID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot update payment method of another user",
		})
	}

	// Update only provided fields
	if req.CardHolderName != "" {
		existing.CardHolderName = req.CardHolderName
	}
	if req.BillingAddress != "" {
		existing.BillingAddress = req.BillingAddress
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.State != "" {
		existing.State = req.State
	}
	if req.PostalCode != "" {
		existing.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		existing.Country = req.Country
	}
	if req.PhoneNumber != "" {
		existing.PhoneNumber = req.PhoneNumber
	}
	existing.IsDefault = req.IsDefault

	result, err := h.paymentMethodRepo.UpdatePaymentMethod(ctx, objectID, existing)
	if err != nil {
		fmt.Printf("[ERROR] Failed to update payment method: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update payment method",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment method updated successfully",
		"data":    result,
	})
}

// DeletePaymentMethod deletes a payment method
func (h *AppHandler) DeletePaymentMethod(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment method ID is required",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment method ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Verify ownership
	existing, err := h.paymentMethodRepo.GetPaymentMethodByID(ctx, objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment method not found",
		})
	}

	if existing.UserID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete payment method of another user",
		})
	}

	if err := h.paymentMethodRepo.DeletePaymentMethod(ctx, objectID); err != nil {
		fmt.Printf("[ERROR] Failed to delete payment method: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete payment method",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment method deleted successfully",
	})
}

// SetDefaultPaymentMethod sets a payment method as default
func (h *AppHandler) SetDefaultPaymentMethod(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment method ID is required",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment method ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Verify ownership
	existing, err := h.paymentMethodRepo.GetPaymentMethodByID(ctx, objectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment method not found",
		})
	}

	if existing.UserID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot set payment method of another user as default",
		})
	}

	if err := h.paymentMethodRepo.SetDefaultPaymentMethod(ctx, objectID, userID.(string)); err != nil {
		fmt.Printf("[ERROR] Failed to set default payment method: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set default payment method",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Default payment method set successfully",
	})
}

// GetSecurityPreferences retrieves security preferences for the logged-in user
func (h *AppHandler) GetSecurityPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch security preferences",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Security preferences retrieved successfully",
		"preferences": user.SecurityPreferences,
	})
}

// SaveSecurityPreferences saves security preferences for the logged-in user
func (h *AppHandler) SaveSecurityPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var preferences appmodels.SecurityPreferences
	if err := c.BodyParser(&preferences); err != nil {
		fmt.Printf("[ERROR] Failed to parse security preferences: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Get user
	user, err := h.userRepo.FindByID(ctx, userID.(string))
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save security preferences",
		})
	}

	// Update security preferences
	user.SecurityPreferences = preferences

	// Save updated user
	if err := h.userRepo.Update(ctx, user.ID, user); err != nil {
		fmt.Printf("[ERROR] Failed to update user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save security preferences",
		})
	}

	fmt.Println("[SUCCESS] Security preferences saved successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Security preferences saved successfully",
		"preferences": user.SecurityPreferences,
	})
}

// GetBanners retrieves all homepage banners for the app
func (h *AppHandler) GetBanners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	travelExperiences, err := h.travelRepo.ListAll(ctx)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch travel experiences: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch banners",
		})
	}

	if travelExperiences == nil {
		travelExperiences = []*adminmodels.TravelExperience{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Banners retrieved successfully",
		"data":    travelExperiences,
	})
}

// GetHomepageBanners retrieves all homepage banner images (ads) for the app
func (h *AppHandler) GetHomepageBanners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	banners, err := h.bannerRepo.ListAll(ctx)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch homepage banners: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch homepage banners",
		})
	}

	if banners == nil {
		banners = []*adminmodels.Banner{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Homepage banners retrieved successfully",
		"data":    banners,
	})
}

// GetSearchBanners retrieves all search results page banner images for the app
func (h *AppHandler) GetSearchBanners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	banners, err := h.searchBannerRepo.ListAll(ctx)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch search banners: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch search banners",
		})
	}

	if banners == nil {
		banners = []*adminmodels.SearchBanner{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Search banners retrieved successfully",
		"data":    banners,
	})
}

// GetPopularLocations retrieves all popular locations for the app
func (h *AppHandler) GetPopularLocations(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	locations, err := h.popularRepo.ListAll(ctx)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch popular locations: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch popular locations",
		})
	}

	if locations == nil {
		locations = []*adminmodels.PopularLocation{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Popular locations retrieved successfully",
		"data":    locations,
	})
}

// CreateHotelBooking creates a new hotel booking with email notification
func (h *AppHandler) CreateHotelBooking(c *fiber.Ctx) error {
	var booking appmodels.HotelBooking

	// Parse request body
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userIDStr := userID.(string)
	fmt.Printf("[BOOKING] Creating hotel booking for user: %s\n", userIDStr)

	// Set booking details
	booking.ID = uuid.New().String()
	booking.UserID = userIDStr
	booking.BookingID = fmt.Sprintf("HBK-%s", uuid.New().String()[:8])

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Create booking in database
	if err := h.hotelBookingRepo.Create(ctx, &booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send email notification (non-blocking)
	// Create a copy of the booking to avoid race conditions in the goroutine
	bookingCopy := booking
	go func() {
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer notifyCancel()

		fmt.Printf("[BOOKING] Sending notification for booking %s, user: %s\n", bookingCopy.BookingID, bookingCopy.UserID)
		if err := h.notificationHelper.SendHotelBookingEmail(notifyCtx, &bookingCopy); err != nil {
			fmt.Printf("[BOOKING] Warning: Failed to send booking notification: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Hotel booking created successfully",
		"booking": booking,
	})
}

// CreateFlightBooking creates a new flight booking with email notification
func (h *AppHandler) CreateFlightBooking(c *fiber.Ctx) error {
	var booking appmodels.FlightBooking

	// Parse request body
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set booking details
	booking.ID = uuid.New().String()
	booking.UserID = userID.(string)
	booking.BookingID = fmt.Sprintf("FBK-%s", uuid.New().String()[:8])

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Create booking in database
	if err := h.flightBookingRepo.Create(ctx, &booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send email notification (non-blocking)
	go func() {
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer notifyCancel()

		if err := h.notificationHelper.SendFlightBookingEmail(notifyCtx, &booking); err != nil {
			fmt.Printf("[BOOKING] Warning: Failed to send booking notification: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Flight booking created successfully",
		"booking": booking,
	})
}

// CreateCarBooking creates a new car booking with email notification
func (h *AppHandler) CreateCarBooking(c *fiber.Ctx) error {
	var booking appmodels.CarBooking

	// Parse request body
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set booking details
	booking.ID = uuid.New().String()
	booking.UserID = userID.(string)
	booking.BookingID = fmt.Sprintf("CBK-%s", uuid.New().String()[:8])

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Create booking in database
	if err := h.carBookingRepo.Create(ctx, &booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send email notification (non-blocking)
	go func() {
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer notifyCancel()

		if err := h.notificationHelper.SendCarBookingEmail(notifyCtx, &booking); err != nil {
			fmt.Printf("[BOOKING] Warning: Failed to send booking notification: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Car booking created successfully",
		"booking": booking,
	})
}

// CreateTransferBooking creates a new transfer booking with email notification
func (h *AppHandler) CreateTransferBooking(c *fiber.Ctx) error {
	var booking appmodels.TransferBooking

	// Parse request body
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set booking details
	booking.ID = uuid.New().String()
	booking.UserID = userID.(string)
	booking.BookingID = fmt.Sprintf("TBK-%s", uuid.New().String()[:8])

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Create booking in database
	if err := h.transferBookingRepo.Create(ctx, &booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send email notification (non-blocking)
	go func() {
		notifyCtx, notifyCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer notifyCancel()

		if err := h.notificationHelper.SendTransferBookingEmail(notifyCtx, &booking); err != nil {
			fmt.Printf("[BOOKING] Warning: Failed to send booking notification: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transfer booking created successfully",
		"booking": booking,
	})
}
