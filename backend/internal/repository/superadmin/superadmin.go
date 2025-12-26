package superadmin

import (
	"context"

	"Anytime_Trip/backend/internal/models/superadmin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SystemConfigRepository handles database operations for system configuration
type SystemConfigRepository struct {
	collection *mongo.Collection
}

// NewSystemConfigRepository creates a new system config repository
func NewSystemConfigRepository(db *mongo.Database) *SystemConfigRepository {
	return &SystemConfigRepository{
		collection: db.Collection("system_config"),
	}
}

// FindByKey finds a config by key
func (r *SystemConfigRepository) FindByKey(ctx context.Context, key string) (*superadmin.SystemConfig, error) {
	var config superadmin.SystemConfig
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateByKey updates a config by key
func (r *SystemConfigRepository) UpdateByKey(ctx context.Context, key string, config *superadmin.SystemConfig) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"key": key},
		bson.M{"$set": config},
	)
	return err
}

// FindAll finds all system configurations
func (r *SystemConfigRepository) FindAll(ctx context.Context) ([]superadmin.SystemConfig, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []superadmin.SystemConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}
