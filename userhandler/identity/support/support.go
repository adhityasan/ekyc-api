package support

import (
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Support struct for modeling Support in mongo collection
type Support struct {
	ID            primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	IdentityID    primitive.ObjectID  `json:"identity_id,omitempty" bson:"identity_id,omitempty"`
	NoHp          string              `schema:"NOMOR_HANDPHONE,omitempty" bson:"nomor_handphone,omitempty" json:"NOMOR_HANDPHONE,omitempty"`
	NoNPWP        string              `schema:"NOMOR_NPWP,omitempty" bson:"nomor_npwp,omitempty" json:"NOMOR_NPWP,omitempty"`
	FotoKTP       *photos.PhotoStruct `schema:"FOTO_KTP" bson:"foto_ktp,omitempty" json:"FOTO_KTP,omitempty"`
	FotoSelfie    *photos.PhotoStruct `schema:"FOTO_SELFIE" bson:"foto_selfie,omitempty" json:"FOTO_SELFIE,omitempty"`
	FotoSelfieKTP *photos.PhotoStruct `schema:"FOTO_SELFIE_KTP" bson:"foto_selfie_ktp,omitempty" json:"FOTO_SELFIE_KTP,omitempty"`
	PasfotoKTP    *photos.PhotoStruct `schema:"PASFOTO_KTP,omitempty" bson:"pasfoto_ktp,omitempty" json:"PASFOTO_KTP,omitempty"`
}
