package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string    `json:"username" bson:"username" binding:"required"`
	Email     string    `json:"email" bson:"email" binding:"required,email"`
	Password  string    `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
