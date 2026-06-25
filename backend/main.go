package main

import (
	"context"
	"log"
	"os"
	"time"

	"travel/backend/config"
	"travel/backend/core/utils"
	"travel/backend/internal/database"
	adminhandlers "travel/backend/internal/handlers/admin"
	apphandlers "travel/backend/internal/handlers/app"
	superadminhandlers "travel/backend/internal/handlers/superadmin"
	"travel/backend/internal/middleware"
	adminmodels "travel/backend/internal/models/admin"
	adminrepo "travel/backend/internal/repository/admin"
	apprepo "travel/backend/internal/repository/app"
	superadminrepo "travel/backend/internal/repository/superadmin"
	adminRoutes "travel/backend/internal/routes/admin"
	appRoutes "travel/backend/internal/routes/app"
	superAdminRoutes "travel/backend/internal/routes/superadmin"
	"travel/backend/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbConn, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	appRepository := apprepo.NewUserRepository(dbConn.DB)
	otpRepository := apprepo.NewOTPRepository(dbConn.DB)
	adminRepository := adminrepo.NewAdminRepository(dbConn.DB)
	passwordResetRepository := adminrepo.NewPasswordResetRepository(dbConn.DB)
	notificationPrefsRepository := adminrepo.NewNotificationPreferencesRepository(dbConn.DB)
	loginAttemptRepository := adminrepo.NewLoginAttemptRepository(dbConn.DB)
	superAdminRepository := superadminrepo.NewSystemConfigRepository(dbConn.DB)
	carBookingRepository := apprepo.NewCarBookingRepository(dbConn.DB)
	flightBookingRepository := apprepo.NewFlightBookingRepository(dbConn.DB)
	hotelBookingRepository := apprepo.NewHotelBookingRepository(dbConn.DB)
	transferBookingRepository := apprepo.NewTransferBookingRepository(dbConn.DB)
	// Support ticket repository
	supportTicketRepository := apprepo.NewSupportTicketRepository(dbConn.DB)
	// Payments repository
	paymentRepository := apprepo.NewPaymentRepository(dbConn.DB)
	// Payment methods repository
	paymentMethodRepository := apprepo.NewPaymentMethodRepository(dbConn.DB)
	flightRepository := adminrepo.NewFlightRepository(dbConn.DB)
	carRepository := adminrepo.NewCarRepository(dbConn.DB)
	hotelRepository := adminrepo.NewHotelRepository(dbConn.DB)
	transferRepository := adminrepo.NewTransferRepository(dbConn.DB)
	bannerRepository := adminrepo.NewBannerRepository(dbConn.DB)
	travelRepository := adminrepo.NewTravelRepository(dbConn.DB)
	popularRepository := adminrepo.NewPopularRepository(dbConn.DB)
	predefinedAnswerRepository := superadminrepo.NewPredefinedAnswerRepository(dbConn.DB)

	if err := ensureDefaultAdmin(adminRepository); err != nil {
		log.Printf("warning: unable to seed default admin user: %v", err)
	}
	if err := ensureDefaultSuperAdmin(adminRepository); err != nil {
		log.Printf("warning: unable to seed default super admin user: %v", err)
	}

	// Initialize WebSocket hub
	chatHub := ws.NewHub()
	go chatHub.Run()
	// Admin notification hub (global) for ticket unlock / notify events
	notifyHub := ws.NewAdminHub()
	go notifyHub.Run()

	// Initialize handlers
	appHandler := apphandlers.NewAppHandler(appRepository, otpRepository, paymentMethodRepository, cfg.JWTSecret)
	adminHandler := adminhandlers.NewAdminHandler(adminRepository, notificationPrefsRepository, passwordResetRepository, appRepository, carBookingRepository, flightBookingRepository, hotelBookingRepository, transferBookingRepository, supportTicketRepository, paymentRepository, flightRepository, carRepository, hotelRepository, transferRepository, bannerRepository, travelRepository, popularRepository, predefinedAnswerRepository, loginAttemptRepository, cfg.JWTSecret, chatHub, notifyHub)
	superAdminHandler := superadminhandlers.NewSuperAdminHandler(superAdminRepository, predefinedAnswerRepository)

	// Initialize Fiber app
	fiberApp := fiber.New()

	// Configure CORS
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(cfg.Environment),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,HX-Request,HX-Target,HX-Current-URL,HX-Trigger",
		AllowCredentials: true,
		ExposeHeaders:    "HX-Redirect,HX-Trigger,HX-Retarget,HX-Reswap",
		MaxAge:           3600,
	}))

	// Root route - redirect to admin login
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/login")
	})

	// Setup static files (relative to backend working dir)
	fiberApp.Static("/static", "./static")
	fiberApp.Static("/", "./templates")
	// Serve uploaded files
	fiberApp.Static("/uploads", "./static/uploads")
	// Serve admin assets (css/images) under /assets
	fiberApp.Static("/assets", "../frontend/admin/src/components/assets")

	// Setup API routes
	api := fiberApp.Group("/api")

	// App routes
	appGroup := api.Group("/app")
	appRoutes.SetupRoutes(appGroup, appHandler, cfg.JWTSecret)

	// Admin routes (with middleware) served under /admin
	adminGroup := fiberApp.Group("/admin")
	// Apply auth middleware only to protected routes, excluding login pages
	adminRoutes.SetupRoutes(adminGroup, adminHandler, cfg.JWTSecret)

	// Super Admin routes (with middleware)
	superAdminGroup := fiberApp.Group("/superadmin")
	superAdminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	superAdminGroup.Use(middleware.SuperAdminMiddleware())
	superAdminRoutes.SetupRoutes(superAdminGroup, superAdminHandler)

	log.Printf("Server starting on :%s...", cfg.Port)
	if err := fiberApp.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// ensureDefaultAdmin creates a default admin user if none exists so login works out of the box.
func ensureDefaultAdmin(repo *adminrepo.AdminRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := repo.FindByEmail(ctx, "admin@travel.app"); err == nil {
		return nil // already exists
	}

	hashed, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	adminUser := &adminmodels.AdminUser{
		ID:          1,
		Username:    "admin",
		Password:    hashed,
		Email:       "admin@travel.app", // stored lowercase
		Role:        "admin",
		Permissions: []string{"all"},
		CreatedAt:   time.Now(),
	}

	return repo.Create(ctx, adminUser)
}

// ensureDefaultSuperAdmin creates a default super admin user if none exists.
func ensureDefaultSuperAdmin(repo *adminrepo.AdminRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := repo.FindByEmail(ctx, "superadmin@travel.app"); err == nil {
		return nil // already exists
	}

	hashed, err := utils.HashPassword("superadmin123")
	if err != nil {
		return err
	}

	superAdminUser := &adminmodels.AdminUser{
		ID:          2,
		Username:    "superadmin",
		Password:    hashed,
		Email:       "superadmin@travel.app", // stored lowercase
		Role:        "superadmin",
		Permissions: []string{"all"},
		CreatedAt:   time.Now(),
	}

	return repo.Create(ctx, superAdminUser)
}

// getAllowedOrigins returns CORS allowed origins based on environment.
// In production, set ALLOWED_ORIGINS env var to a comma-separated list of origins.
func getAllowedOrigins(environment string) string {
	if allowed := os.Getenv("ALLOWED_ORIGINS"); allowed != "" {
		return allowed
	}
	if environment == "production" {
		return "*"
	}
	return "http://localhost:3000,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:5173,http://127.0.0.1:8080"
}
