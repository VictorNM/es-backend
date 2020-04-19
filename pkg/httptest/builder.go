package httptest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type RequestBuilder struct {
	req *http.Request

	body map[string]interface{}
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		body: make(map[string]interface{}),
	}
}

func (b *RequestBuilder) GET(target string) *RequestBuilder {
	b.req = httptest.NewRequest(http.MethodGet, target, nil)
	return b
}

func (b *RequestBuilder) POST(target string) *RequestBuilder {
	b.req = httptest.NewRequest(http.MethodGet, target, nil)
	return b
}

func (b *RequestBuilder) PATCH(target string) *RequestBuilder {
	b.req = httptest.NewRequest(http.MethodGet, target, nil)
	return b
}

func (b *RequestBuilder) PUT(target string) *RequestBuilder {
	b.req = httptest.NewRequest(http.MethodGet, target, nil)
	return b
}

func (b *RequestBuilder) DELETE(target string) *RequestBuilder {
	b.req = httptest.NewRequest(http.MethodGet, target, nil)
	return b
}

func (b *RequestBuilder) AddHeader(key, value string) *RequestBuilder {
	b.req.Header.Add(key, value)
	return b
}

func (b *RequestBuilder) AddJSON(o interface{}) *RequestBuilder {
	data, err := json.Marshal(o)
	if err != nil {
		log.Panicf("marshal JSON error %v", err)
	}

	var m map[string]interface{}

	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Panicf("unmarshal JSON error %v", err)
	}

	for k, v := range m {
		b.body[k] = v
	}

	return b
}

func (b *RequestBuilder) AddJSONField(key string, value interface{}) *RequestBuilder {
	b.body[key] = value
	return b
}

func (b *RequestBuilder) Build() *http.Request {
	body, err := json.Marshal(b.body)
	if err != nil {
		log.Panicf("marshal JSON error %v", err)
	}

	b.req.Body = ioutil.NopCloser(bytes.NewReader(body))

	return b.req
}
