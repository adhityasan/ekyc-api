package status

import "go.mongodb.org/mongo-driver/bson/primitive"

// Status as struct model fot Stat collection
type Status struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IdentityID     primitive.ObjectID `json:"identity_id,omitempty" bson:"identity_id,omitempty"`
	FaceComparison float64
	DataPercentage int
}
