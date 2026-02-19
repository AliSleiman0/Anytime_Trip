package app

import (
	"context"
	"time"

	appmodels "Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatbotRatingRepository provides DB access for chatbot ratings
type ChatbotRatingRepository struct {
	collection *mongo.Collection
}

// NewChatbotRatingRepository creates a new repository
func NewChatbotRatingRepository(db *mongo.Database) *ChatbotRatingRepository {
	return &ChatbotRatingRepository{collection: db.Collection("chatbot_ratings")}
}

// Create creates a new chatbot rating
func (r *ChatbotRatingRepository) Create(ctx context.Context, rating *appmodels.ChatbotRating) (primitive.ObjectID, error) {
	rating.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, rating)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetByUserID gets all ratings by a specific user
func (r *ChatbotRatingRepository) GetByUserID(ctx context.Context, userID string) ([]appmodels.ChatbotRating, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []appmodels.ChatbotRating
	if err := cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

// GetAll gets all chatbot ratings
func (r *ChatbotRatingRepository) GetAll(ctx context.Context) ([]appmodels.ChatbotRating, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ratings []appmodels.ChatbotRating
	if err := cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

// GetAverageRating calculates the average rating
func (r *ChatbotRatingRepository) GetAverageRating(ctx context.Context) (float64, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":     nil,
			"average": bson.M{"$avg": "$rating"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Average float64 `bson:"average"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Average, nil
}
