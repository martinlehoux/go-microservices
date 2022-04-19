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

func WriteResponse(w http.ResponseWriter, code int, data Data) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	PanicOnError(err)
}

type Data = map[string]interface{}
