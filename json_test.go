package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTools_ReadJSON(t *testing.T) {
	var tests = []struct {
		name               string
		json               string
		errorExpected      bool
		maxSize            int
		allowUnknownFields bool
	}{
		{"Valid JSON", `{"foo":"bar"}`, false, 1024, true},
		{"Invalid JSON", "Hello, World!", true, 1024, false},
		{"Invalid JSON - missing values", `{"foo":}`, true, 1024, false},
		{"Invalid JSON - incorrect type", `{"foo": 1}`, true, 1024, false},
		{"More than one JSON", `{"foo":"bar"}{"hello":"world"}`, true, 1024, false},
		{"Empty JSON", ``, true, 1024, false},
		{"Syntax error", `{"foo":"bar"`, true, 1024, false},
		{"Unknown field", `{"hello":"world"}`, true, 1024, false},
		{"Allow unknown fields", `{"hello":"world"}`, false, 1024, true},
		{"Missing field name", `{foo:"bar"}`, true, 1024, true},
		{"File is too large", `{"foo":"bar"}`, true, 1, false},
	}
	var tools Tools
	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Set the max file size
			tools.MaxJSONSize = entry.maxSize

			// Allow or disllow unknown fields
			tools.AllowUnknownFields = entry.allowUnknownFields

			// Declare a variable that will read the decoded JSON
			var decodedJSON struct {
				Foo string `json:"foo"`
			}

			// Create a request with a body
			// Cast the requests body to a slice of bytes
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(entry.json)))
			if err != nil {
				t.Log("Error:", err)
			}
			// Create a response recorder
			resp := httptest.NewRecorder()

			// Call the ReadJSON function
			err = tools.ReadJSON(resp, req, &decodedJSON)

			if entry.errorExpected && err == nil {
				t.Error("expected an error, but received none")
			}

			if !entry.errorExpected && err != nil {
				t.Errorf("expected no error, but received %+v", err)
			}

			// Close the request body to avoid resource leaks
			req.Body.Close()
		})
	}
}

func TestTools_WriteJSON(t *testing.T) {
	tests := []struct {
		hdr   string
		value string
	}{
		{"FOO", "BAR"},
		{"HELLO", "WORLD"},
	}
	for _, entry := range tests {
		t.Run(entry.hdr, func(t *testing.T) {

			var tools Tools
			// Create a new response recorder
			resp := httptest.NewRecorder()
			JSONpayload := JSONResponse{
				Error:   false,
				Message: "foo",
			}

			headers := make(http.Header)
			headers.Add(entry.hdr, entry.value)

			err := tools.WriteJSON(resp, http.StatusOK, JSONpayload, headers)
			if err != nil {
				t.Errorf("failed to write JSON: %+v", err)
			}

			// Check headers
			setHeader := resp.Header().Get(entry.hdr)

			if setHeader != entry.value {
				t.Errorf("expected to receive header %s, but received %s", entry.hdr, setHeader)
			}
		})
	}

}

func TestTools_ErrorJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"Not acceptable", http.StatusNotAcceptable},
		{"Forbidden", http.StatusForbidden},
		{"Gateway timeout", http.StatusGatewayTimeout},
	}
	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			var tools Tools
			// Create response recorder
			resp := httptest.NewRecorder()
			err := tools.ErrorJSON(resp, errors.New("error"), entry.statusCode)
			if err != nil {
				t.Error(err)
			}

			// Check the response
			var JSONPayload JSONResponse
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&JSONPayload)

			if err != nil {
				t.Error("received error when decoding JSON:", err)
			}

			if !JSONPayload.Error {
				t.Error("error set to false to JSON, but it should be set to true")
			}

			if resp.Code != entry.statusCode {
				t.Errorf("expected status code %d, but received %d", entry.statusCode, resp.Code)
			}
		})
	}

}
