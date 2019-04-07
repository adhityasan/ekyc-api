package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhityasan/ekyc-api/userhandler/identity"

	"github.com/adhityasan/ekyc-api/imagehandler"
	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"
	"github.com/adhityasan/ekyc-api/userhandler/ocr"
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

// Ocr request handler for /ocr route
func Ocr(response http.ResponseWriter, request *http.Request) {
	// Sementara pakai buffer, next pakai Pii untuk return objectID
	imgChan := make(chan interface{})
	imageStruct, _, err := photos.PhotoStructHandler("OCR_IMAGE", request)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(writeResponseByte(err.Error(), nil))
		return
	}

	adapter := &imagehandler.AwsAdapter{}
	go adapter.Read(imageStruct.Data, imgChan)
	ocrRes := <-imgChan

	var ocrreq ocr.Request
	ocrreq.GenerateToken()
	ocrreq.OcrImage = imageStruct
	ocrreq.OcrResult = ocrRes
	errsave := ocrreq.Save()
	if errsave != nil {
		log.Println(errsave)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(writeResponseByte(errsave.Error(), nil))
		return
	}

	var customData ocr.CustomResponse
	customData.ID = ocrreq.ID
	customData.ClientID = ocrreq.ClientID
	customData.OcrResult = ocrreq.OcrResult

	response.Header().Set("Ocrtoken", ocrreq.Token)
	response.Write(writeResponseByte("", customData))
}

// Register to assign new fake data to Pii collection
func Register(response http.ResponseWriter, request *http.Request) {

	userIdentity, _ := identity.DecodeFormPost(request)
	isExistLoc, _ := userIdentity.Exist()
	var message string

	if isExistLoc {
		errgrep := userIdentity.GrepData()
		message = "KYC databases"
		if errgrep != nil {
			log.Println(errgrep)
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(writeResponseByte(errgrep.Error(), userIdentity.Nik))
			return
		}
	} else {
		errgrep := userIdentity.GrepDataFromDukcapil()
		message = "Dukcapil databases"
		if errgrep != nil {
			log.Println(errgrep)
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(writeResponseByte(errgrep.Error(), userIdentity.Nik))
			return
		}
	}

	response.Write(writeResponseByte(message, userIdentity))
}
