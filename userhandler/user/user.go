package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct for modeling User in mongo collection
type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email      string
	Password   string
	IdentityID primitive.ObjectID `json:"identity_id,omitempty" bson:"identity_id,omitempty"`
}
