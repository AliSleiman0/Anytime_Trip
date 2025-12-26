package main

import (
	"log"

	"Anytime_Trip/backend/config"
	"Anytime_Trip/backend/internal/database"
	adminhandlers "Anytime_Trip/backend/internal/handlers/admin"
	apphandlers "Anytime_Trip/backend/internal/handlers/app"
	superadminhandlers "Anytime_Trip/backend/internal/handlers/superadmin"
	"Anytime_Trip/backend/internal/middleware"
	adminrepo "Anytime_Trip/backend/internal/repository/admin"
	apprepo "Anytime_Trip/backend/internal/repository/app"
	superadminrepo "Anytime_Trip/backend/internal/repository/superadmin"
	adminRoutes "Anytime_Trip/backend/internal/routes/admin"
	appRoutes "Anytime_Trip/backend/internal/routes/app"
	superAdminRoutes "Anytime_Trip/backend/internal/routes/superadmin"

	"github.com/gofiber/fiber/v2"
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
	adminRepository := adminrepo.NewAdminRepository(dbConn.DB)
	superAdminRepository := superadminrepo.NewSystemConfigRepository(dbConn.DB)

	// Initialize handlers
	appHandler := apphandlers.NewAppHandler(appRepository)
	adminHandler := adminhandlers.NewAdminHandler(adminRepository)
	superAdminHandler := superadminhandlers.NewSuperAdminHandler(superAdminRepository)

	// Initialize Fiber app
	fiberApp := fiber.New()

	// Setup static files
	fiberApp.Static("/static", "../../frontend")
	fiberApp.Static("/", "../../templates")

	// Setup API routes
	api := fiberApp.Group("/api")

	// App routes
	appGroup := api.Group("/app")
	appRoutes.SetupRoutes(appGroup, appHandler)

	// Admin routes (with middleware)
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware)
	adminGroup.Use(middleware.AdminMiddleware)
	adminRoutes.SetupRoutes(adminGroup, adminHandler)

	// Super Admin routes (with middleware)
	superAdminGroup := api.Group("/superadmin")
	superAdminGroup.Use(middleware.AuthMiddleware)
	superAdminGroup.Use(middleware.SuperAdminMiddleware)
	superAdminRoutes.SetupRoutes(superAdminGroup, superAdminHandler)

	log.Printf("Server starting on :%s...", cfg.Port)
	if err := fiberApp.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
