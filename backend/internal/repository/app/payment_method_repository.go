package app

import (
	"context"
	"fmt"
	"time"

	appmodels "travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaymentMethodRepository provides DB access for payment methods
type PaymentMethodRepository struct {
	collection *mongo.Collection
}

// NewPaymentMethodRepository creates a new repository
func NewPaymentMethodRepository(db *mongo.Database) *PaymentMethodRepository {
	return &PaymentMethodRepository{collection: db.Collection("payment_methods")}
}

// AddPaymentMethod adds a new payment method for a user
func (r *PaymentMethodRepository) AddPaymentMethod(ctx context.Context, pm *appmodels.PaymentMethod) (*appmodels.PaymentMethod, error) {
	pm.ID = primitive.NewObjectID()
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = pm.CreatedAt

	// If this is set as default, unset all other defaults for this user
	if pm.IsDefault {
		_, err := r.collection.UpdateMany(
			ctx,
			bson.M{"user_id": pm.UserID},
			bson.M{"$set": bson.M{"is_default": false}},
		)
		if err != nil {
			fmt.Printf("[ERROR] Failed to unset other defaults: %v\n", err)
		}
	}

	_, err := r.collection.InsertOne(ctx, pm)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[SUCCESS] Payment method added successfully for user %s\n", pm.UserID)
	return pm, nil
}

// GetPaymentMethodsByUserID retrieves all payment methods for a user
func (r *PaymentMethodRepository) GetPaymentMethodsByUserID(ctx context.Context, userID string) ([]appmodels.PaymentMethod, error) {
	opts := options.Find().SetSort(bson.D{{Key: "is_default", Value: -1}, {Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID, "is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []appmodels.PaymentMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}

	return methods, nil
}

// GetPaymentMethodByID retrieves a specific payment method
func (r *PaymentMethodRepository) GetPaymentMethodByID(ctx context.Context, id primitive.ObjectID) (*appmodels.PaymentMethod, error) {
	var method appmodels.PaymentMethod
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&method)
	if err != nil {
		return nil, err
	}
	return &method, nil
}

// UpdatePaymentMethod updates a payment method
func (r *PaymentMethodRepository) UpdatePaymentMethod(ctx context.Context, id primitive.ObjectID, pm *appmodels.PaymentMethod) (*appmodels.PaymentMethod, error) {
	pm.UpdatedAt = time.Now()

	// If setting as default, unset other defaults for this user
	if pm.IsDefault {
		_, err := r.collection.UpdateMany(
			ctx,
			bson.M{"user_id": pm.UserID, "_id": bson.M{"$ne": id}},
			bson.M{"$set": bson.M{"is_default": false}},
		)
		if err != nil {
			fmt.Printf("[ERROR] Failed to unset other defaults: %v\n", err)
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": pm},
		opts,
	).Decode(&pm)

	if err != nil {
		return nil, err
	}

	fmt.Printf("[SUCCESS] Payment method updated successfully\n")
	return pm, nil
}

// DeletePaymentMethod soft deletes a payment method
func (r *PaymentMethodRepository) DeletePaymentMethod(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("[SUCCESS] Payment method deleted successfully\n")
	return nil
}

// SetDefaultPaymentMethod sets a payment method as default for a user
func (r *PaymentMethodRepository) SetDefaultPaymentMethod(ctx context.Context, id primitive.ObjectID, userID string) error {
	// Unset all other defaults
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "_id": bson.M{"$ne": id}},
		bson.M{"$set": bson.M{"is_default": false}},
	)
	if err != nil {
		return err
	}

	// Set this one as default
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_default": true, "updated_at": time.Now()}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("[SUCCESS] Default payment method set successfully\n")
	return nil
}
