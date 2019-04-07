package photos

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// PhotoStruct standart Struct for image
type PhotoStruct struct {
	Data []byte `schema:"data,omitempty" bson:"data,omitempty"`
	Name string `schema:"name,omitempty" bson:"name,omitempty"`
	Size int64  `schema:"size,omitempty" bson:"size,omitempty"`
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

	return newPhotoStruct, file, nil
}
