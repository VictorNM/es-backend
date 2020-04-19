package httptest

import (
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

func assertStatusEqual(t *testing.T, w *httptest.ResponseRecorder, code int) {
	assert.Equal(t, code, w.Code)
}
