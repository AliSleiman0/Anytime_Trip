package admin

import "time"

// TravelExperience represents a homepage travel experience item
type TravelExperience struct {
	ID        string    `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	ImagePath string    `bson:"image_path" json:"image_path"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
