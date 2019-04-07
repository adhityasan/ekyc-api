package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/adhityasan/ekyc-api/config"
	"github.com/adhityasan/ekyc-api/db"
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbhost = config.Of.Mongo.Host
var dbport = config.Of.Mongo.Port
var dburl = config.Of.Mongo.URL
var dbname = config.Of.DBModules["identity"].Db
var dbcoll = config.Of.DBModules["identity"].Coll

// Identity struct for modeling Identity in mongo collection
type Identity struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"ID,omitempty"`
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
	Foto             *photos.PhotoStruct `schema:"FOTO,omitempty" bson:"foto,omitempty" json:"-"`
	SupportID        primitive.ObjectID  `schema:"SUPPORT_ID,omitempty" bson:"support_id,omitempty" json:"-"`
}

// Validate validate current identity data
func (identity *Identity) Validate() bool {
	nikmatch, _ := regexp.MatchString("^(11|12|13|14|15|16|17|18|19|21|31|32|33|34|35|36|51|52|53|61|62|63|64|65|71|72|73|74|75|76|81|82|91|92)(0[^0]|[^0][0-9]){2}([04][^0]|[1256][0-9]|[37][01])(0[1-9]|1[012])[0-9]{2}([0-9]{0,3}[^0][0-9]{0,3})$", identity.Nik)
	return nikmatch
}

// Save Save identity to mongo dataidentity collection
func (identity *Identity) Save() error {

	exist, _ := identity.Exist()
	if exist {
		return errors.New("Identity data exist, Identity.ID has been set")
	}

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

// Exist Check Pii data existance in local database
func (identity *Identity) Exist() (bool, error) {
	_, cancel, _, collection, errconn := db.OpenConnection(10, dburl, dbname, dbcoll)
	if errconn != nil {
		return false, errconn
	}

	decodepoint := new(Identity)
	errfind := collection.FindOne(context.TODO(), bson.M{"nik": identity.Nik}).Decode(decodepoint)
	defer cancel()

	if errfind != nil {
		return false, errfind
	}

	pointerID := &identity.ID
	*pointerID = decodepoint.ID

	return true, nil
}

// GrepData grep all current identity data by its ID From Local Database
func (identity *Identity) GrepData() error {
	_, cancel, _, collection, errconn := db.OpenConnection(10, dburl, dbname, dbcoll)
	if errconn != nil {
		return errconn
	}

	errfind := collection.FindOne(context.TODO(), bson.M{"_id": identity.ID}).Decode(&identity)
	if errfind != nil {
		return errfind
	}
	defer cancel()

	return nil
}

// GrepDataFromDukcapil  grep all current identity data by its ID From Dukcapil
func (identity *Identity) GrepDataFromDukcapil() error {
	requstBody, err := json.Marshal(map[string]string{
		"NIK": identity.Nik,
	})

	req, err := http.NewRequest("POST", config.Of.Dukcapil.Endpoint, bytes.NewBuffer(requstBody))
	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		log.Println("errDo", errDo)
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Fail generate fake identity. Bad request")
	}

	defer resp.Body.Close()

	var decoded struct {
		Content []interface{}
	}

	errDecode := json.NewDecoder(resp.Body).Decode(&decoded)
	if errDecode != nil {
		log.Println("errdecode", errDecode)
		return errDecode
	}

	thecontents, errMarshalContent := json.Marshal(decoded.Content[0])
	if errMarshalContent != nil {
		log.Println("errMarshalContent", errMarshalContent)
		return errMarshalContent
	}

	json.Unmarshal(thecontents, &identity)

	return nil
}
