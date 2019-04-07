package ocr

import (
	"crypto/rand"
	"fmt"

	"github.com/adhityasan/ekyc-api/config"
	"github.com/adhityasan/ekyc-api/db"
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbhost = config.Of.Mongo.Host
var dbport = config.Of.Mongo.Port
var dburl = config.Of.Mongo.URL
var dbname = config.Of.DBModules["ocr"].Db
var dbcoll = config.Of.DBModules["ocr"].Coll

// Request struct for modeling Request in mongo collection
type Request struct {
	ID        primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID  `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Token     string              `json:"token,omitempty" bson:"token,omitempty"`
	ClientID  primitive.ObjectID  `json:"client_id,omitempty" bson:"client_id,omitempty"`
	OcrImage  *photos.PhotoStruct `json:"ocr_image,omitmpety" bson:"ocr_image,omitempty"`
	OcrResult interface{}         `json:"ocr_result,omitempty" bson:"ocr_result,omitempty"`
}

// CustomResponse for handle response data from the controller
type CustomResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID  primitive.ObjectID `json:"client_id,omitempty" bson:"client_id,omitempty"`
	OcrResult interface{}        `json:"ocr_result,omitempty" bson:"ocr_result,omitempty"`
}

// GenerateToken genereate random string into RequestTOken
func (r *Request) GenerateToken() {
	b := make([]byte, 12)
	rand.Read(b)
	token := &r.Token
	*token = fmt.Sprintf("%x", b)
}

// Save save Request struct into datarequest collection
func (r *Request) Save() error {
	ctx, cancel, _, collection, err := db.OpenConnection(10, dburl, dbname, dbcoll)

	res, err := collection.InsertOne(ctx, r)
	defer cancel()
	if err != nil {
		return err
	}

	newid := &r.ID
	*newid = res.InsertedID.(primitive.ObjectID)

	return nil
}
