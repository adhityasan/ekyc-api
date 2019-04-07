package unittest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhityasan/ekyc-api/userhandler/identity/assigner"
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"

	"github.com/adhityasan/ekyc-api/userhandler/identity"
)

type controllerResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeResponseByte(msg string, data interface{}) []byte {
	var resp controllerResponse
	if len(msg) == 0 {
		msg = "Success"
	}
	resp.Message = msg
	resp.Data = data
	res, _ := json.Marshal(resp)

	return res
}

// AssignFakeIdentity to assign new fake data to Identity collection
func AssignFakeIdentity(response http.ResponseWriter, request *http.Request) {

	var userIdentity identity.Identity
	var message string
	contentType := request.Header.Get("Content-Type")

	switch contentType {
	case "multipart/form-data":

		message = "Data assigned from requests formData"
		userIdentity, _ := identity.DecodeFormPost(request)

		isExistLoc, _ := userIdentity.Exist()
		if isExistLoc {
			userIdentity.GrepData()
			message = "data based on NIK:" + userIdentity.Nik + " already exists"
			response.WriteHeader(http.StatusBadRequest)
		}

	case "application/json":

		type reqjsonstruct struct {
			Nik  string `json:"NIK,omitempty"`
			Foto string `json:"FOTO,omitempty"`
		}

		var formjson reqjsonstruct
		errDecode := json.NewDecoder(request.Body).Decode(&formjson)
		if errDecode != nil {
			log.Println(errDecode)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(writeResponseByte(errDecode.Error(), formjson.Nik))
			return
		}

		errGenerateFromDukcapil := assigner.DukcapilSimulatorAssigner(formjson.Nik, &userIdentity)
		if errGenerateFromDukcapil != nil {
			log.Println(errGenerateFromDukcapil)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(writeResponseByte(errGenerateFromDukcapil.Error(), formjson.Nik))
			return
		}

		var fotostruct photos.PhotoStruct
		errGeneratePhotoStruct := fotostruct.GenerateStructFromURL(formjson.Foto)
		if errGeneratePhotoStruct != nil {
			log.Println(errGeneratePhotoStruct)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(writeResponseByte(errGeneratePhotoStruct.Error(), formjson.Nik))
			return
		}
		userIdentity.Foto = &fotostruct

	default:

		response.WriteHeader(http.StatusBadRequest)
		message = "Content-Type that are allowed are only application/json and multipart/form-data"

	}

	errSave := userIdentity.Save()
	if errSave != nil {
		log.Println(errSave)
		response.Write(writeResponseByte(errSave.Error(), userIdentity))
		return
	}

	response.Write(writeResponseByte(message, userIdentity))
}

// GrepData grep Identity Data by user NIK
func GrepData(response http.ResponseWriter, request *http.Request) {

	userIdentity, _ := identity.DecodeFormPost(request)
	isExistLoc, _ := userIdentity.Exist()
	var message string

	if isExistLoc {
		errgrep := userIdentity.GrepData()
		message = "Data exist in KYC databases"
		if errgrep != nil {
			log.Println(errgrep)
			response.Write(writeResponseByte(errgrep.Error(), userIdentity.Nik))
			return
		}
	} else {
		errgrep := userIdentity.GrepDataFromDukcapil()
		message = "Dukcapil databases"
		if errgrep != nil {
			log.Println(errgrep)
			response.Write(writeResponseByte(errgrep.Error(), userIdentity.Nik))
			return
		}
	}

	response.Write(writeResponseByte(message, userIdentity))
}
