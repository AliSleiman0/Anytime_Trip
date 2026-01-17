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

// SupportTicketRepository provides DB access for support tickets
type SupportTicketRepository struct {
	collection *mongo.Collection
}

// NewSupportTicketRepository creates a new repository
func NewSupportTicketRepository(db *mongo.Database) *SupportTicketRepository {
	return &SupportTicketRepository{collection: db.Collection("support_tickets")}
}

// GetAllTickets returns all tickets sorted by created_at desc
func (r *SupportTicketRepository) GetAllTickets(ctx context.Context) ([]appmodels.SupportTicket, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var out []appmodels.SupportTicket
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTicketByID finds a ticket by ObjectID
func (r *SupportTicketRepository) GetTicketByID(ctx context.Context, id primitive.ObjectID) (*appmodels.SupportTicket, error) {
	var t appmodels.SupportTicket
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTicketByTicketID finds a ticket by its short ticket_id field
func (r *SupportTicketRepository) GetTicketByTicketID(ctx context.Context, ticketID string) (*appmodels.SupportTicket, error) {
	var t appmodels.SupportTicket
	if err := r.collection.FindOne(ctx, bson.M{"ticket_id": ticketID}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

// AddReply appends a reply to a ticket and updates updated_at/status
func (r *SupportTicketRepository) AddReply(ctx context.Context, ticketID primitive.ObjectID, reply appmodels.TicketReply) error {
	reply.ID = primitive.NewObjectID()
	reply.CreatedAt = time.Now()
	update := bson.M{
		"$push": bson.M{"replies": reply},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": ticketID}, update)
	return err
}

// UpdateTicketStatus sets the ticket status and updates updated_at
func (r *SupportTicketRepository) UpdateTicketStatus(ctx context.Context, ticketID primitive.ObjectID, status string) error {
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": ticketID}, update)
	return err
}

// UpdateTicketPriority sets priority and updates updated_at
func (r *SupportTicketRepository) UpdateTicketPriority(ctx context.Context, ticketID primitive.ObjectID, priority string) error {
	update := bson.M{"$set": bson.M{"priority": priority, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": ticketID}, update)
	return err
}
