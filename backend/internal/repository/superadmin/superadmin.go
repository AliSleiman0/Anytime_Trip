package superadmin

import (
	"context"
	"time"

	"travel/backend/internal/models/superadmin"

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

// PredefinedAnswerRepository handles database operations for predefined answers
type PredefinedAnswerRepository struct {
	collection *mongo.Collection
}

// NewPredefinedAnswerRepository creates a new predefined answer repository
func NewPredefinedAnswerRepository(db *mongo.Database) *PredefinedAnswerRepository {
	return &PredefinedAnswerRepository{
		collection: db.Collection("predefined_answers"),
	}
}

// Create creates a new predefined answer
func (r *PredefinedAnswerRepository) Create(ctx context.Context, answer *superadmin.PredefinedAnswer) error {
	answer.CreatedAt = time.Now()
	answer.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, answer)
	return err
}

// FindAll finds all predefined answers
func (r *PredefinedAnswerRepository) FindAll(ctx context.Context) ([]superadmin.PredefinedAnswer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var answers []superadmin.PredefinedAnswer
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil
}

// FindByID finds a predefined answer by ID
func (r *PredefinedAnswerRepository) FindByID(ctx context.Context, id int) (*superadmin.PredefinedAnswer, error) {
	var answer superadmin.PredefinedAnswer
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&answer)
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

// Update updates a predefined answer
func (r *PredefinedAnswerRepository) Update(ctx context.Context, id int, answer *superadmin.PredefinedAnswer) error {
	answer.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{
			"shortcut":   answer.Shortcut,
			"answer":     answer.Answer,
			"updated_at": answer.UpdatedAt,
		}},
	)
	return err
}

// Delete deletes a predefined answer by ID
func (r *PredefinedAnswerRepository) Delete(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
