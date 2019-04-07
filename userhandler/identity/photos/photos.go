package photos

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// PhotoStruct standart Struct for image
type PhotoStruct struct {
	Data   []byte               `schema:"data,omitempty" bson:"data,omitempty"`
	Name   string               `schema:"name,omitempty" bson:"name,omitempty"`
	Size   int64                `schema:"size,omitempty" bson:"size,omitempty"`
	Header textproto.MIMEHeader `schema:"header,omitempty" bson:"header,omitempty"`
}

// PhotoStructHandler create image Struct with Data , Name, Size, Header from a request
func PhotoStructHandler(fieldname string, r *http.Request) (*PhotoStruct, multipart.File, error) {
	file, handler, err := r.FormFile(fieldname)

	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	newPhotoStruct := new(PhotoStruct)
	newPhotoStruct.Data = fileBytes
	newPhotoStruct.Name = handler.Filename
	newPhotoStruct.Size = handler.Size
	newPhotoStruct.Header = handler.Header

	return newPhotoStruct, file, nil
}
