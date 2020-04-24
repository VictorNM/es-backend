package api

import (
	"encoding/json"
	"github.com/victornm/es-backend/api/internal/jsontest"
	"log"
	"net/http"
	"net/http/httptest"
)

// newRequest marshal the body before pass to httptest.NewRequest
func newRequest(method, target string, body interface{}) *http.Request {
	return jsontest.NewRequest(method, target, body)
}

func signIn(handler http.Handler, email, password string) string {
	req := jsontest.WrapPOST("/api/users/sign-in", nil).
		SetBasicAuth(email, password).
		Unwrap()

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
	log.Println(w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &res)
	if err != nil {
		log.Panicf("invalid json %v", err)
	}

	return res
}
