package admin

import (
	"context"
	"errors"
	"strings"

	"Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AdminRepository handles database operations for admin users
type AdminRepository struct {
	collection *mongo.Collection
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{
		collection: db.Collection("admins"),
	}
}

// FindByEmail finds an admin by email
func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*admin.AdminUser, error) {
	var adminUser admin.AdminUser
	normalized := strings.ToLower(email)
	err := r.collection.FindOne(ctx, bson.M{"email": normalized}).Decode(&adminUser)
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// FindByID finds an admin by ID
func (r *AdminRepository) FindByID(ctx context.Context, id string) (*admin.AdminUser, error) {
	var adminUser admin.AdminUser
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&adminUser)
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// FindAll finds all admin users
func (r *AdminRepository) FindAll(ctx context.Context) ([]admin.AdminUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var admins []admin.AdminUser
	if err := cursor.All(ctx, &admins); err != nil {
		return nil, err
	}
	return admins, nil
}

// Create creates a new admin user
func (r *AdminRepository) Create(ctx context.Context, adminUser *admin.AdminUser) error {
	_, err := r.collection.InsertOne(ctx, adminUser)
	return err
}

// UpdatePassword updates an admin's password after verifying the current password
func (r *AdminRepository) UpdatePassword(ctx context.Context, emailOrID string, currentPassword string, newPassword string) error {
	// Try to find the admin by email first (since JWT stores email in Subject)
	var adminUser admin.AdminUser
	err := r.collection.FindOne(ctx, bson.M{"email": strings.ToLower(emailOrID)}).Decode(&adminUser)
	if err != nil {
		// If not found by email, try by _id
		err = r.collection.FindOne(ctx, bson.M{"_id": emailOrID}).Decode(&adminUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("admin not found")
			}
			return err
		}
	}

	// Verify current password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password using email as identifier (more reliable than _id)
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"email": adminUser.Email},
		bson.M{"$set": bson.M{"password": string(hashedPassword)}},
	)

	return err
}

// UpdateByEmail updates an admin user by email
func (r *AdminRepository) UpdateByEmail(ctx context.Context, email string, adminUser *admin.AdminUser) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"email": strings.ToLower(email)},
		bson.M{"$set": adminUser},
	)
	return err
}

// UpdateCurrency updates an admin's currency preference
func (r *AdminRepository) UpdateCurrency(ctx context.Context, emailOrID string, currency string) error {
	// Try to find the admin by email first (since JWT stores email in Subject)
	var adminUser admin.AdminUser
	err := r.collection.FindOne(ctx, bson.M{"email": strings.ToLower(emailOrID)}).Decode(&adminUser)
	if err != nil {
		// If not found by email, try by _id
		err = r.collection.FindOne(ctx, bson.M{"_id": emailOrID}).Decode(&adminUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("admin not found")
			}
			return err
		}
	}

	// Update currency using email as identifier
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"email": adminUser.Email},
		bson.M{"$set": bson.M{"currency": currency}},
	)

	return err
}

// UpdateLanguage updates an admin's language preference
func (r *AdminRepository) UpdateLanguage(ctx context.Context, emailOrID string, language string) error {
	// Try to find the admin by email first (since JWT stores email in Subject)
	var adminUser admin.AdminUser
	err := r.collection.FindOne(ctx, bson.M{"email": strings.ToLower(emailOrID)}).Decode(&adminUser)
	if err != nil {
		// If not found by email, try by _id
		err = r.collection.FindOne(ctx, bson.M{"_id": emailOrID}).Decode(&adminUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("admin not found")
			}
			return err
		}
	}

	// Update language using email as identifier
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"email": adminUser.Email},
		bson.M{"$set": bson.M{"language": language}},
	)

	return err
}
