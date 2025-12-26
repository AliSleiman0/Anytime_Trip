package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds the database connection
type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewDB creates a new MongoDB connection
func NewDB(uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &DB{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if db.Client != nil {
		return db.Client.Disconnect(ctx)
	}
	return nil
}
