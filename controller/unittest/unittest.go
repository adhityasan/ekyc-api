package unittest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhityasan/ekyc-api/userhandler/identity"
	"github.com/adhityasan/ekyc-api/userhandler/identity/assigner"
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

		json.NewDecoder(request.Body).Decode(&userIdentity)
		err := assigner.DukcapilSimulatorAssigner(userIdentity.Nik, &userIdentity)
		if err != nil {
			log.Println(err)
			response.Write(writeResponseByte(err.Error(), userIdentity.Nik))
			return
		}

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
