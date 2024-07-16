package web

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	internalServerErrorMessage = "Internal Server Error"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func WriteJSON(rw http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func InternalServerError(rw http.ResponseWriter, err error) {
	errorMessage := ErrorMessage{Message: internalServerErrorMessage}
	writeError := WriteJSON(rw, http.StatusInternalServerError, errorMessage)
	if writeError != nil {
		log.Println("Error writing error message: ", writeError.Error())
	}
	log.Print("Internal error server: ", err.Error())
}

func MakeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(data)
		_, _ = w.Write(body)
	}
}
