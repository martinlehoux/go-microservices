package common

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ID = uuid.UUID

func CreateID() ID {
	return uuid.New()
}

func PanicOnError(err error) {
	if err != nil {
		log.Panicf("unhandled error: %s", err.Error())
	}
}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	PanicOnError(err)
}

func WriteError(w http.ResponseWriter, code int, err error) {
	log.Printf("error: %s", err.Error())
	WriteResponse(w, code, ErrorDto{Error: err.Error()})
}

type ErrorDto struct {
	Error string `json:"error"`
}

type OperationDto struct {
	Success bool `json:"success"`
}

type AnyDto map[string]interface{}
