package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/core/utils"
	"Anytime_Travel/backend/internal/database"
	adminhandlers "Anytime_Travel/backend/internal/handlers/admin"
	apphandlers "Anytime_Travel/backend/internal/handlers/app"
	superadminhandlers "Anytime_Travel/backend/internal/handlers/superadmin"
	"Anytime_Travel/backend/internal/middleware"
	adminmodels "Anytime_Travel/backend/internal/models/admin"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	apprepo "Anytime_Travel/backend/internal/repository/app"
	superadminrepo "Anytime_Travel/backend/internal/repository/superadmin"
	adminRoutes "Anytime_Travel/backend/internal/routes/admin"
	appRoutes "Anytime_Travel/backend/internal/routes/app"
	superAdminRoutes "Anytime_Travel/backend/internal/routes/superadmin"
	"Anytime_Travel/backend/internal/ws"

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
	// Support ticket repository
	supportTicketRepository := apprepo.NewSupportTicketRepository(dbConn.DB)
	// Payments repository
	paymentRepository := apprepo.NewPaymentRepository(dbConn.DB)
	flightRepository := adminrepo.NewFlightRepository(dbConn.DB)
	carRepository := adminrepo.NewCarRepository(dbConn.DB)
	hotelRepository := adminrepo.NewHotelRepository(dbConn.DB)
	bannerRepository := adminrepo.NewBannerRepository(dbConn.DB)
	travelRepository := adminrepo.NewTravelRepository(dbConn.DB)
	popularRepository := adminrepo.NewPopularRepository(dbConn.DB)
	predefinedAnswerRepository := superadminrepo.NewPredefinedAnswerRepository(dbConn.DB)

	if err := ensureDefaultAdmin(adminRepository); err != nil {
		log.Printf("warning: unable to seed default admin user: %v", err)
	}

	// Initialize WebSocket hub
	chatHub := ws.NewHub()
	go chatHub.Run()
	// Admin notification hub (global) for ticket unlock / notify events
	notifyHub := ws.NewAdminHub()
	go notifyHub.Run()

	// Initialize handlers
	appHandler := apphandlers.NewAppHandler(appRepository, otpRepository, cfg.JWTSecret)
	adminHandler := adminhandlers.NewAdminHandler(adminRepository, notificationPrefsRepository, passwordResetRepository, appRepository, carBookingRepository, flightBookingRepository, hotelBookingRepository, supportTicketRepository, paymentRepository, flightRepository, carRepository, hotelRepository, bannerRepository, travelRepository, popularRepository, predefinedAnswerRepository, loginAttemptRepository, cfg.JWTSecret, chatHub, notifyHub)
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
	// Serve admin assets (css/images) under /assets
	fiberApp.Static("/assets", "../frontend/admin/src/components/assets")

	// Setup API routes
	api := fiberApp.Group("/api")

	// App routes
	appGroup := api.Group("/app")
	appRoutes.SetupRoutes(appGroup, appHandler)

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

	if _, err := repo.FindByEmail(ctx, "admin@anytime.com"); err == nil {
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
		Email:       "admin@anytime.com", // stored lowercase
		Role:        "admin",
		Permissions: []string{"all"},
		CreatedAt:   time.Now(),
	}

	return repo.Create(ctx, adminUser)
}

// getAllowedOrigins returns CORS allowed origins based on environment
func getAllowedOrigins(environment string) string {
	if environment == "production" {
		// In production, specify exact domains
		return "https://yourdomain.com,https://www.yourdomain.com,https://admin.yourdomain.com"
	}
	// In development, allow localhost origins
	return "http://localhost:3000,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:5173,http://127.0.0.1:8080"
}
