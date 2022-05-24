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
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	PanicOnError(err)
}

// Depracated: use WriteErrors instead
func WriteError(w http.ResponseWriter, code int, err error) {
	WriteResponse(w, code, ErrorDto{Error: err.Error()})
}

func WriteErrors(w http.ResponseWriter, code int, errors []error) {
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

type Middleware func(http.Handler) http.Handler

func CommonMiddlewareConstructor(logger Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.Header().Set("Content-Type", "application/json")

			ctx := logger.Init(req.Context())
			requestId := uuid.New().String()
			logger.With(ctx, "requestId", requestId).With(ctx, "method", req.Method).With(ctx, "url", req.URL.String()).Info(ctx, "request started")

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
