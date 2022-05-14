package common

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ID = uuid.UUID

func CreateID() ID {
	return uuid.New()
}

func ParseID(id string) (ID, error) {
	parsedId, err := uuid.Parse(id)
	return parsedId, err
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

// Depracated: use WriteErrors instead
func WriteError(w http.ResponseWriter, code int, err error) {
	log.Printf("error: %s", err.Error())
	WriteResponse(w, code, ErrorDto{Error: err.Error()})
}

func WriteErrors(w http.ResponseWriter, code int, errors []error) {
	log.Printf("errors: %s", errors)
	var dto ErrorsDto
	for _, err := range errors {
		dto.Errors = append(dto.Errors, err.Error())
	}
	WriteResponse(w, code, dto)
}

// Depracated: use ErrorsDto instead
type ErrorDto struct {
	Error string `json:"error"`
}

type ErrorsDto struct {
	Errors []string `json:"errors"`
}

type OperationDto struct {
	Success bool `json:"success"`
}

type AnyDto map[string]interface{}

var ErrURLNotFound = errors.New("url not found")
