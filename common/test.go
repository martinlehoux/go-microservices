package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func PrepareRequest(method string, target string, payload interface{}) http.Request {
	body, err := json.Marshal(payload)
	PanicOnError(err)

	req := httptest.NewRequest(method, target, bytes.NewReader(body))

	req.Header.Add("Content-Type", "application/json")
	return *req
}
