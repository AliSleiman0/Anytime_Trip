package admin

import (
	"context"

	"Anytime_Trip/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
