package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// newRequest marshal the body before pass to httptest.NewRequest
func newRequest(method, target string, body interface{}) *http.Request {
	if reader, ok := body.(io.Reader); ok {
		return httptest.NewRequest(method, target, reader)
	}

	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(method, target, bytes.NewReader(b))
	return req
}

func signIn(handler http.Handler, email, password string) string {
	req := newRequest(http.MethodPost, "/api/users/sign-in", nil)
	req.SetBasicAuth(email, password)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	return getDataAsMap(w)["token"].(string)
}

func getErrors(w *httptest.ResponseRecorder) []Error {
	return getResponse(w).Errors
}

func getData(w *httptest.ResponseRecorder) interface{} {
	return getResponse(w).Data
}

func getDataAsMap(w *httptest.ResponseRecorder) map[string]interface{} {
	data, ok := getData(w).(map[string]interface{})
	if !ok {
		log.Panicf("response data is not a %T", map[string]interface{}{})
	}

	return data
}

func getResponse(w *httptest.ResponseRecorder) *BaseResponse {
	var res *BaseResponse
	err := json.Unmarshal([]byte(w.Body.String()), &res)
	if err != nil {
		log.Panicf("invalid json %v", err)
	}

	return res
}
