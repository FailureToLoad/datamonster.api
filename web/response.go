package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func WriteJSON(rw http.ResponseWriter, status int, data any) error {
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

func InternalServerError(rw http.ResponseWriter, errMessage string) {
	writeError := WriteJSON(rw, http.StatusInternalServerError, errMessage)
	if writeError != nil {
		log.Println("Error writing error message: ", writeError.Error())
	}
	log.Print("Internal error server: ", errMessage)
}

type unauthorized struct {
	Access string  `json:"access"`
	Reason *string `json:"message,omitempty"`
}

func Unauthorized(rw http.ResponseWriter, message *string) {
	content := unauthorized{
		Access: "unauthorized",
	}
	if message != nil {
		content.Reason = message
	}
	writeError := WriteJSON(rw, http.StatusUnauthorized, content)
	if writeError != nil {
		log.Println("Error writing error message: ", writeError.Error())
	}
	log.Println("Unauthorized: ", message)
}

func MakeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(data)
		_, _ = w.Write(body)
	}
}
