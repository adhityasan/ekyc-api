package ocr

import (
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request struct for modeling Request in mongo collection
type Request struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RequestToken string
	ClientID     primitive.ObjectID `json:"client_id,omitempty" bson:"client_id,omitempty"`
	OcrImage     *photos.PhotoStruct
	OcrResult    []byte
}
