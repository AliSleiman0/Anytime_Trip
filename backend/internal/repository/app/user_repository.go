package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository handles database operations for app users
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*app.User, error) {
	var user app.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *app.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, id string, user *app.User) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// UpdateLastLogin sets the last_login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string, at time.Time) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"last_login": at}})
	return err
}

// CountTotalUsers returns the total number of users
func (r *UserRepository) CountTotalUsers(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	return count, err
}

// CountNewUsersToday returns the number of users created today
func (r *UserRepository) CountNewUsersToday(ctx context.Context) (int64, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	})
	return count, err
}

// CountActiveUsers returns the number of users created in the last 30 days
func (r *UserRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": thirtyDaysAgo,
		},
	})
	return count, err
}

// FindAll returns users with pagination support
func (r *UserRepository) FindAll(ctx context.Context, limit int64, skip int64) ([]app.User, error) {
	findOptions := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cur, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []app.User
	for cur.Next(ctx) {
		var u app.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateEmail updates the email for a user
func (r *UserRepository) UpdateEmail(ctx context.Context, id string, email string) error {
	update := bson.M{
		"$set": bson.M{
			"email": email,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
