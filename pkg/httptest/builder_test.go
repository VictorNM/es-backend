package httptest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestBuilder_GET(t *testing.T) {
	b := NewRequestBuilder()
	req := b.GET("/api/ping").Build()

	assert.Equal(t, "/api/ping", req.URL.String())
}

func TestRequestBuilder_AddHeader(t *testing.T) {
	b := NewRequestBuilder()
	req := b.GET("/api/ping").
		AddHeader("Authorization", "Bearer 123").
		Build()

	assert.Equal(t, "Bearer 123", req.Header.Get("Authorization"))
}

func TestRequestBuilder_POST(t *testing.T) {
	b := NewRequestBuilder()
	req := b.POST("/api/ping").
		AddJSON(map[string]interface{}{
			"email":    "foo@bar.com",
			"password": "1234abcd",
		}).
		AddJSONField("username", "victor").
		Build()

	var m map[string]interface{}
	_ = json.NewDecoder(req.Body).Decode(&m)
	assert.Equal(t, m["email"], "foo@bar.com")
	assert.Equal(t, m["password"], "1234abcd")
	assert.Equal(t, m["username"], "victor")
}
