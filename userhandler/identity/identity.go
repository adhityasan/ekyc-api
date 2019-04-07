package identity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identity struct for modeling Identity in mongo collection
type Identity struct {
	ID               primitive.ObjectID `schema:"_id,omitempty" bson:"_id,omitempty" json:"_id,omitempty"`
	Nik              string             `schema:"NIK,omitempty" bson:"nik,omitempty" json:"NIK,omitempty"`
	EktpStatus       bool               `schema:"EKTP_STATUS,omitempty" bson:"ektp_status,omitempty" json:"EKTP_STATUS,omitempty"`
	NamaLengkap      string             `schema:"NAMA_LENGKAP,omitempty" bson:"nama_lengkap,omitempty" json:"NAMA_LENGKAP,omitempty"`
	TanggalLahir     string             `schema:"TANGGAL_LAHIR,omitempty" bson:"tanggal_lahir,omitempty" json:"TANGGAL_LAHIR,omitempty"`
	TempatLahir      string             `schema:"TEMPAT_LAHIR,omitempty" bson:"tempat_lahir,omitempty" json:"TEMPAT_LAHIR,omitempty"`
	Kewarganegaraan  string             `schema:"KEWARGANEGARAAN,omitempty" bson:"kewarganegaraan,omitempty" json:"KEWARGANEGARAAN,omitempty"`
	PendidikanAkhir  string             `schema:"PENDIDIKAN_AKHIR,omitempty" bson:"pendidikan_akhir,omitempty" json:"PENDIDIKAN_AKHIR,omitempty"`
	NoKK             string             `schema:"NOMOR_KARTU_KELUARGA,omitempty" bson:"nomor_kartu_keluarga,omitempty" json:"NOMOR_KARTU_KELUARGA,omitempty"`
	Alamat           string             `schema:"ALAMAT,omitempty" bson:"alamat,omitempty" json:"ALAMAT,omitempty"`
	Rt               string             `schema:"RT,omitempty" bson:"rt,omitempty" json:"RT,omitempty"`
	Rw               string             `schema:"RW,omitempty" bson:"rw,omitempty" json:"RW,omitempty"`
	Kelurahan        string             `schema:"KELURAHAN,omitempty" bson:"kelurahan,omitempty" json:"KELURAHAN,omitempty"`
	Kecamatan        string             `schema:"KECAMATAN,omitempty" bson:"kecamatan,omitempty" json:"KECAMATAN,omitempty"`
	Kabupaten        string             `schema:"KABUPATEN,omitempty" bson:"kabupaten,omitempty" json:"KABUPATEN,omitempty"`
	Provinsi         string             `schema:"PROVINSI,omitempty" bson:"provinsi,omitempty" json:"PROVINSI,omitempty"`
	Agama            string             `schema:"AGAMA,omitempty" bson:"agama,omitempty" json:"AGAMA,omitempty"`
	Pekerjaan        string             `schema:"PEKERJAAN,omitempty" bson:"pekerjaan,omitempty" json:"PEKERJAAN,omitempty"`
	JenisKelamin     string             `schema:"JENIS_KELAMIN,omitempty" bson:"jenis_kelamin,omitempty" json:"JENIS_KELAMIN,omitempty"`
	StatusPerkawinan string             `schema:"STATUS_PERKAWINAN,omitempty" bson:"status_perkawinan,omitempty" json:"STATUS_PERKAWINAN,omitempty"`
	SupportID        primitive.ObjectID `schema:"SUPPORT_ID,omitempty" bson:"support_id,omitempty" json:"SUPPORT_ID,omitempty"`
}
