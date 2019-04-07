package identity

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adhityasan/ekyc-api/config"
	"github.com/adhityasan/ekyc-api/db"
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbhost = config.Of.Mongo.Host
var dbport = config.Of.Mongo.Port
var dburl = config.Of.Mongo.URL
var dbname = config.Of.DBModules["identity"].Db
var dbcoll = config.Of.DBModules["identity"].Coll

// Identity struct for modeling Identity in mongo collection
type Identity struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	Nik              string              `schema:"NIK,omitempty" bson:"nik,omitempty" json:"NIK,omitempty"`
	EktpStatus       bool                `schema:"EKTP_STATUS,omitempty" bson:"ektp_status,omitempty" json:"EKTP_STATUS,omitempty"`
	NamaLengkap      string              `schema:"NAMA_LENGKAP,omitempty" bson:"nama_lengkap,omitempty" json:"NAMA_LENGKAP,omitempty"`
	TanggalLahir     string              `schema:"TANGGAL_LAHIR,omitempty" bson:"tanggal_lahir,omitempty" json:"TANGGAL_LAHIR,omitempty"`
	TempatLahir      string              `schema:"TEMPAT_LAHIR,omitempty" bson:"tempat_lahir,omitempty" json:"TEMPAT_LAHIR,omitempty"`
	Kewarganegaraan  string              `schema:"KEWARGANEGARAAN,omitempty" bson:"kewarganegaraan,omitempty" json:"KEWARGANEGARAAN,omitempty"`
	PendidikanAkhir  string              `schema:"PENDIDIKAN_AKHIR,omitempty" bson:"pendidikan_akhir,omitempty" json:"PENDIDIKAN_AKHIR,omitempty"`
	NoKK             string              `schema:"NOMOR_KARTU_KELUARGA,omitempty" bson:"nomor_kartu_keluarga,omitempty" json:"NOMOR_KARTU_KELUARGA,omitempty"`
	Alamat           string              `schema:"ALAMAT,omitempty" bson:"alamat,omitempty" json:"ALAMAT,omitempty"`
	Rt               string              `schema:"RT,omitempty" bson:"rt,omitempty" json:"RT,omitempty"`
	Rw               string              `schema:"RW,omitempty" bson:"rw,omitempty" json:"RW,omitempty"`
	Kelurahan        string              `schema:"KELURAHAN,omitempty" bson:"kelurahan,omitempty" json:"KELURAHAN,omitempty"`
	Kecamatan        string              `schema:"KECAMATAN,omitempty" bson:"kecamatan,omitempty" json:"KECAMATAN,omitempty"`
	Kabupaten        string              `schema:"KABUPATEN,omitempty" bson:"kabupaten,omitempty" json:"KABUPATEN,omitempty"`
	Provinsi         string              `schema:"PROVINSI,omitempty" bson:"provinsi,omitempty" json:"PROVINSI,omitempty"`
	Agama            string              `schema:"AGAMA,omitempty" bson:"agama,omitempty" json:"AGAMA,omitempty"`
	Pekerjaan        string              `schema:"PEKERJAAN,omitempty" bson:"pekerjaan,omitempty" json:"PEKERJAAN,omitempty"`
	JenisKelamin     string              `schema:"JENIS_KELAMIN,omitempty" bson:"jenis_kelamin,omitempty" json:"JENIS_KELAMIN,omitempty"`
	StatusPerkawinan string              `schema:"STATUS_PERKAWINAN,omitempty" bson:"status_perkawinan,omitempty" json:"STATUS_PERKAWINAN,omitempty"`
	Foto             *photos.PhotoStruct `schema:"FOTO,omitempty" bson:"foto,omitempty" json:"FOTO,omitempty"`
	SupportID        primitive.ObjectID  `schema:"SUPPORT_ID,omitempty" bson:"support_id,omitempty" json:"SUPPORT_ID,omitempty"`
}

// Save Save identity to mongo dataidentity collection
func (identity *Identity) Save() error {
	ctx, cancel, _, collection, err := db.OpenConnection(10, dburl, dbname, dbcoll)

	res, err := collection.InsertOne(ctx, identity)
	defer cancel()
	if err != nil {
		return err
	}

	newid := &identity.ID
	*newid = res.InsertedID.(primitive.ObjectID)

	return nil
}

// DecodeFormPost decode the formPost data in requests form-data and assign it to Pii Struct
func DecodeFormPost(request *http.Request) (*Identity, error) {

	request.ParseMultipartForm(10 << 20)

	fd := request.PostForm
	newIdentity := new(Identity)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if fd.Get("TANGGAL_LAHIR") != "" {
		parsedtgllahir, errParse := time.Parse("2006-01-02", fd.Get("TANGGAL_LAHIR"))
		if errParse != nil {
			return nil, errors.New("Cannot parse tanggallahir decode")
		}
		tgllahirString := parsedtgllahir.String()
		fd.Set("TANGGAL_LAHIR", tgllahirString)
	}

	err := decoder.Decode(newIdentity, fd)

	if err != nil {
		errdetail := fmt.Sprintf("Fail to decode request form-data to new Identity data Struct : %s\n", err)
		return nil, errors.New(errdetail)
	}

	newIdentity.Foto, _, err = photos.PhotoStructHandler("FOTO", request)

	return newIdentity, nil
}
