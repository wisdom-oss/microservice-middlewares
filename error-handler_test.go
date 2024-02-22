package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
)

var errorMap = map[string]wisdomType.WISdoMError{
	"WISDOM_TEST": {
		Type:   "WISDOM_TEST",
		Status: 400,
		Title:  "Bad Request",
		Detail: "This is only a test",
	},
}

var nativeError = errors.New("this is a native error")

func TestErrorHandler_NativeError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ErrorHandler)

	r.Get("/", func(writer http.ResponseWriter, r *http.Request) {
		errorChannel := r.Context().Value(ErrorChannelName).(chan<- interface{})
		statusChannel := r.Context().Value(StatusChannelName).(<-chan bool)

		errorChannel <- nativeError
		<-statusChannel
	})

	r.ServeHTTP(recorder, request)
	res := recorder.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("response status is incorrect, got %d, expected, %d", res.StatusCode, http.StatusInternalServerError)
	}

	if res.Header.Get("Content-Type") != "application/problem+json; charset=utf-8" {
		t.Errorf("response content type is incorrect, got '%s', expected '%s'", res.Header.Get("Content-Type"), "application/problem+json")
	}

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("response could not be decoded: %s", err.Error())
	}

	_, validType := response["type"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'type' field to be 'string', got '%T'", response["type"])
	}

	_, validType = response["status"].(float64)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'status' field to be 'int', got '%T'", response["status"])
	}

	_, validType = response["title"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'title' field to be 'string', got '%T'", response["title"])
	}

	_, validType = response["detail"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'detail' field to be 'string', got '%T'", response["detail"])
	}

	_, validType = response["instance"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'instance' field to be 'string', got '%T'", response["instance"])
	}

	_, validType = response["error"].(string)
	if !validType {
		t.Errorf("rfc9457 extension violated, expected 'error' field to be 'string', got '%T'", response["error"])
	}

	if hostname, _ := os.Hostname(); response["instance"].(string) != hostname {
		t.Errorf("instance is not hostname, got '%s', expected '%s'", response["instance"].(string), hostname)
	}
}

func TestErrorHandler_WISdoMError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ErrorHandler)

	r.Get("/", func(writer http.ResponseWriter, r *http.Request) {
		errorChannel := r.Context().Value(ErrorChannelName).(chan<- interface{})
		statusChannel := r.Context().Value(StatusChannelName).(<-chan bool)

		errorChannel <- errorMap["WISDOM_TEST"]
		<-statusChannel
	})

	r.ServeHTTP(recorder, request)
	res := recorder.Result()

	if res.StatusCode != errorMap["WISDOM_TEST"].Status {
		t.Errorf("response is incorrect, got %d, expected %d", res.StatusCode, errorMap["WISDOM_TEST"].Status)
	}

	if res.Header.Get("Content-Type") != "application/problem+json; charset=utf-8" {
		t.Errorf("response content type is incorrect, got '%s', expected '%s'", res.Header.Get("Content-Type"), "application/problem+json")
	}

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("response could not be decoded: %s", err.Error())
	}

	_, validType := response["type"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'type' field to be 'string', got '%T'", response["type"])
	}

	_, validType = response["status"].(float64)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'status' field to be 'int', got '%T'", response["status"])
	}

	_, validType = response["title"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'title' field to be 'string', got '%T'", response["title"])
	}

	_, validType = response["detail"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'detail' field to be 'string', got '%T'", response["detail"])
	}

	_, validType = response["instance"].(string)
	if !validType {
		t.Errorf("rfc9457 violated, expected 'instance' field to be 'string', got '%T'", response["instance"])
	}

	if hostname, _ := os.Hostname(); response["instance"].(string) != hostname {
		t.Errorf("instance is not hostname, got '%s', expected '%s'", response["instance"].(string), hostname)
	}

	if response["type"].(string) != errorMap["WISDOM_TEST"].Type {
		t.Errorf("type field wrong, got '%s', expected '%s'", response["type"].(string), errorMap["WISDOM_TEST"].Type)
	}

	if response["title"].(string) != errorMap["WISDOM_TEST"].Title {
		t.Errorf("title field wrong, got '%s', expected '%s'", response["title"].(string), errorMap["WISDOM_TEST"].Title)
	}

	if response["detail"].(string) != errorMap["WISDOM_TEST"].Detail {
		t.Errorf("detail field wrong, got '%s', expected '%s'", response["detail"].(string), errorMap["WISDOM_TEST"].Detail)
	}
}

func TestErrorHandler_InvalidTypeSupplied(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(ErrorHandler)

	r.Get("/", func(writer http.ResponseWriter, r *http.Request) {
		errorChannel := r.Context().Value(ErrorChannelName).(chan<- interface{})
		statusChannel := r.Context().Value(StatusChannelName).(<-chan bool)

		errorChannel <- "invalid-type"
		<-statusChannel
	})

	r.ServeHTTP(recorder, request)
	res := recorder.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("response status is incorrect, got %d, expected, %d", res.StatusCode, http.StatusInternalServerError)
	}

	if res.Header.Get("Content-Type") != "application/problem+json; charset=utf-8" {
		t.Errorf("response content type is incorrect, got '%s', expected '%s'", res.Header.Get("Content-Type"), "application/problem+json")
	}

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("response could not be decoded: %s", err.Error())
	}

	if hostname, _ := os.Hostname(); response["instance"].(string) != hostname {
		t.Errorf("instance is not hostname, got '%s', expected '%s'", response["instance"].(string), hostname)
	}
}
