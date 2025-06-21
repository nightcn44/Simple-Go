package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Content   string             `json:"content" bson:"content" binding:"required"`
	AuthorID  primitive.ObjectID `json:"author_id" bson:"author_id"` // Reference to User
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}