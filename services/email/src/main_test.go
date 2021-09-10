package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func createMockRequest(headers map[string]string) http.Request {
	req := httptest.NewRequest("GET", "/noop", nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return *req
}

func createMockResponseWriter(statusCode int, headers map[string]string, body string) httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
	return *w
}

func createMockBodyJSON(fields interface{}) []byte {
	body, _ := json.Marshal(fields)
	return body
}

func TestAuthentication(t *testing.T) {
	type test struct {
		env  map[string]string
		r    http.Request
		want bool
	}

	tests := []test{
		{
			// protected, correct key: success
			map[string]string{"API_KEY": "ABC123"},
			createMockRequest(map[string]string{"X-API-KEY": "ABC123"}),
			true,
		},
		{
			// unprotected, no key: success
			map[string]string{},
			createMockRequest(map[string]string{}),
			true,
		},
		{
			// unprotected, no key: success
			map[string]string{"API_KEY": ""},
			createMockRequest(map[string]string{}),
			true,
		},
		{
			// unprotected, empty key: success
			map[string]string{"API_KEY": ""},
			createMockRequest(map[string]string{"X-API-KEY": ""}),
			true,
		},
		{
			// unprotected, with key: success
			map[string]string{"API_KEY": ""},
			createMockRequest(map[string]string{"X-API-KEY": "ABC123"}),
			true,
		},
		{
			// protected, wrong key: failure
			map[string]string{"API_KEY": "ABC123"},
			createMockRequest(map[string]string{"X-API-KEY": "ABC124"}),
			false,
		},
		{
			// protected, no key: failure
			map[string]string{"API_KEY": "ABC123"},
			createMockRequest(map[string]string{}),
			false,
		},
		{
			// protected, empty key: failure
			map[string]string{"API_KEY": "ABC123"},
			createMockRequest(map[string]string{"X-API-KEY": ""}),
			false,
		},
	}

	for _, tc := range tests {
		for key, value := range tc.env {
			os.Setenv(key, value)
		}
		isAuthenticated := authentication(&tc.r)
		if isAuthenticated != tc.want {
			t.Errorf("authentication incorrect: got %v, want %v", isAuthenticated, tc.want)
		}
		for key := range tc.env {
			os.Unsetenv(key)
		}
	}
}

func TestGenerateResponse(t *testing.T) {
	type test struct {
		w          *httptest.ResponseRecorder
		statusCode int
		body       []byte
		want       httptest.ResponseRecorder
	}

	tests := []test{
		{
			httptest.NewRecorder(),
			200,
			createMockBodyJSON(map[string]interface{}{"success": true}),
			createMockResponseWriter(200, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"success\":true}"),
		},
		{
			httptest.NewRecorder(),
			204,
			nil,
			createMockResponseWriter(204, map[string]string{"Content-Type": "application/json; charset=utf-8"}, ""),
		},
		{
			httptest.NewRecorder(),
			400,
			createMockBodyJSON(map[string]interface{}{"error": []string{"error 1", "error 2"}}),
			createMockResponseWriter(400, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"error\":[\"error 1\",\"error 2\"]}"),
		},
		{
			httptest.NewRecorder(),
			500,
			createMockBodyJSON(map[string]interface{}{"error": "Server error."}),
			createMockResponseWriter(500, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"error\":\"Server error.\"}"),
		},
	}

	for _, tc := range tests {
		generateResponse(tc.w, tc.statusCode, tc.body)
		result := tc.w.Result()
		wantResult := tc.want.Result()
		if result.StatusCode != wantResult.StatusCode {
			t.Errorf("generateResponse StatusCode: got %v, want %v", result.StatusCode, wantResult.StatusCode)
		}
		if !reflect.DeepEqual(result.Header, wantResult.Header) {
			t.Errorf("generateResponse Header: got %v, want %v", result.Header, wantResult.Header)
		}
		body, _ := ioutil.ReadAll(result.Body)
		wantBody, _ := ioutil.ReadAll(wantResult.Body)
		if string(body) != string(wantBody) {
			t.Errorf("generateResponse Body: got %v, want %v", string(body), string(wantBody))
		}
	}
}

func TestSuccessResponse(t *testing.T) {
	type test struct {
		w          *httptest.ResponseRecorder
		statusCode int
		fields     interface{}
		want       httptest.ResponseRecorder
	}

	tests := []test{
		{
			httptest.NewRecorder(),
			200,
			map[string]interface{}{"success": true},
			createMockResponseWriter(200, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"success\":true}"),
		},
		{
			httptest.NewRecorder(),
			204,
			nil,
			createMockResponseWriter(204, map[string]string{"Content-Type": "application/json; charset=utf-8"}, ""),
		},
	}

	for _, tc := range tests {
		successResponse(tc.w, tc.statusCode, tc.fields)
		result := tc.w.Result()
		wantResult := tc.want.Result()
		if result.StatusCode != wantResult.StatusCode {
			t.Errorf("successResponse StatusCode: got %v, want %v", result.StatusCode, wantResult.StatusCode)
		}
		if !reflect.DeepEqual(result.Header, wantResult.Header) {
			t.Errorf("successResponse Header: got %v, want %v", result.Header, wantResult.Header)
		}
		body, _ := ioutil.ReadAll(result.Body)
		wantBody, _ := ioutil.ReadAll(wantResult.Body)
		if string(body) != string(wantBody) {
			t.Errorf("successResponse Body: got %v, want %v", string(body), string(wantBody))
		}
	}
}

func TestUserErrorResponse(t *testing.T) {
	type test struct {
		w            *httptest.ResponseRecorder
		statusCode   int
		errorMessage string
		want         httptest.ResponseRecorder
	}

	tests := []test{
		{
			httptest.NewRecorder(),
			400,
			"Bad request.",
			createMockResponseWriter(400, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"error\":\"Bad request.\"}"),
		},
		{
			httptest.NewRecorder(),
			404,
			"Not found.",
			createMockResponseWriter(404, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"error\":\"Not found.\"}"),
		},
	}

	for _, tc := range tests {
		userErrorResponse(tc.w, tc.statusCode, tc.errorMessage)
		result := tc.w.Result()
		wantResult := tc.want.Result()
		if result.StatusCode != wantResult.StatusCode {
			t.Errorf("successResponse StatusCode: got %v, want %v", result.StatusCode, wantResult.StatusCode)
		}
		if !reflect.DeepEqual(result.Header, wantResult.Header) {
			t.Errorf("successResponse Header: got %v, want %v", result.Header, wantResult.Header)
		}
		body, _ := ioutil.ReadAll(result.Body)
		wantBody, _ := ioutil.ReadAll(wantResult.Body)
		if string(body) != string(wantBody) {
			t.Errorf("successResponse Body: got %v, want %v", string(body), string(wantBody))
		}
	}
}

func TestServerErrorResponse(t *testing.T) {
	type test struct {
		w    *httptest.ResponseRecorder
		want httptest.ResponseRecorder
	}

	tests := []test{
		{
			httptest.NewRecorder(),
			createMockResponseWriter(500, map[string]string{"Content-Type": "application/json; charset=utf-8"}, "{\"error\":\"Server error\"}"),
		},
	}

	for _, tc := range tests {
		serverErrorResponse(tc.w)
		result := tc.w.Result()
		wantResult := tc.want.Result()
		if result.StatusCode != wantResult.StatusCode {
			t.Errorf("successResponse StatusCode: got %v, want %v", result.StatusCode, wantResult.StatusCode)
		}
		if !reflect.DeepEqual(result.Header, wantResult.Header) {
			t.Errorf("successResponse Header: got %v, want %v", result.Header, wantResult.Header)
		}
		body, _ := ioutil.ReadAll(result.Body)
		wantBody, _ := ioutil.ReadAll(wantResult.Body)
		if string(body) != string(wantBody) {
			t.Errorf("successResponse Body: got %v, want %v", string(body), string(wantBody))
		}
	}
}
