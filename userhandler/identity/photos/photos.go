package photos

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
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

// GenerateStructFromURL generate PhotoStruct data from image url
func (currphoto *PhotoStruct) GenerateStructFromURL(urlImage string) error {
	filename := path.Base(urlImage)

	resp, errGet := http.Get(urlImage)
	if errGet != nil {
		log.Println("http error get from" + urlImage)
		return errGet
	}
	defer resp.Body.Close()

	pointdir := "/tmp/" + filename
	file, errCreate := os.Create(pointdir)
	if errCreate != nil {
		log.Println("os error create file " + urlImage)
		return errCreate
	}
	defer file.Close()

	_, errCopy := io.Copy(file, resp.Body)
	if errCopy != nil {
		log.Println("io error copy resp.Body to " + urlImage)
		log.Fatal(errCopy)
	}

	fileInfo, _ := file.Stat()
	fileBytes, _ := ioutil.ReadFile(pointdir)

	pointerData := &currphoto.Data
	pointerName := &currphoto.Name
	pointerSize := &currphoto.Size

	*pointerData = fileBytes
	*pointerName = fileInfo.Name()
	*pointerSize = fileInfo.Size()

	return nil
}
