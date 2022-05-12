package common

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type RequestBuilder struct {
	req http.Request
}

func NewRequestBuilder(method string, target string) RequestBuilder {
	req := httptest.NewRequest(method, target, nil)
	return RequestBuilder{req: *req}
}

func (builder RequestBuilder) WithPayload(payload interface{}) RequestBuilder {
	body, err := json.Marshal(payload)
	PanicOnError(err)

	builder.req.Body = ioutil.NopCloser(bytes.NewReader(body))
	builder.req.Header.Add("Content-Type", "application/json")
	return builder
}

func (builder RequestBuilder) WithToken(token Token, privateKey rsa.PrivateKey) RequestBuilder {
	rawToken, _ := SignToken(token, privateKey)
	builder.req.Header.Add("Authorization", "Bearer "+base64.StdEncoding.EncodeToString(rawToken))
	return builder
}

func (builder RequestBuilder) Build() http.Request {
	return builder.req
}
