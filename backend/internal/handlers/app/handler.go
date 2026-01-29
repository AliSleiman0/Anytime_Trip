package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
