//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/core/utils"
	"Anytime_Travel/backend/internal/database"
	adminmodels "Anytime_Travel/backend/internal/models/admin"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
)

func main() {
	log.Println("Starting admin user seeding...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbConn, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize admin repository
	adminRepo := adminrepo.NewAdminRepository(dbConn.DB)

	// Seed admin users
	if err := seedAdminUsers(adminRepo); err != nil {
		log.Fatalf("Failed to seed admin users: %v", err)
	}

	log.Println("✓ Admin user seeding completed successfully!")
	log.Println("\nDefault credentials:")
	log.Println("  Admin: admin@anytime.com / admin123")
	log.Println("  Super Admin: superadmin@anytime.com / super123")
}

// seedAdminUsers seeds initial admin users into the database
func seedAdminUsers(repo *adminrepo.AdminRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	admins := []struct {
		ID          int
		Username    string
		Email       string
		Password    string
		Role        string
		Permissions []string
	}{
		{
			ID:          1,
			Username:    "admin",
			Email:       "admin@anytime.com",
			Password:    "admin123",
			Role:        "admin",
			Permissions: []string{"all"},
		},
		{
			ID:          2,
			Username:    "superadmin",
			Email:       "superadmin@anytime.com",
			Password:    "super123",
			Role:        "superadmin",
			Permissions: []string{"all"},
		},
		{
			ID:          3,
			Username:    "admin1",
			Email:       "tannousszzz16@gmail.com",
			Password:    "anyadmin123",
			Role:        "superadmin",
			Permissions: []string{"all"},
		},
	}

	for _, admin := range admins {
		// Check if admin already exists
		if _, err := repo.FindByEmail(ctx, admin.Email); err == nil {
			log.Printf("Admin user %s already exists, skipping...", admin.Email)
			continue
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(admin.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", admin.Email, err)
			return err
		}

		// Create admin user
		adminUser := &adminmodels.AdminUser{
			ID:          admin.ID,
			Username:    admin.Username,
			Password:    hashedPassword,
			Email:       admin.Email,
			Role:        admin.Role,
			Permissions: admin.Permissions,
			CreatedAt:   time.Now(),
		}

		if err := repo.Create(ctx, adminUser); err != nil {
			log.Printf("Failed to create admin user %s: %v", admin.Email, err)
			return err
		}

		log.Printf("✓ Created admin user: %s (email: %s, password: %s)", admin.Username, admin.Email, admin.Password)
	}

	return nil
}
