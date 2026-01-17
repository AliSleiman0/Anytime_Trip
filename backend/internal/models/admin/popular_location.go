package admin

import "time"

// PopularLocation represents a homepage popular location item
type PopularLocation struct {
	ID        string    `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Location  string    `bson:"location" json:"location"`
	ImagePath string    `bson:"image_path" json:"image_path"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
