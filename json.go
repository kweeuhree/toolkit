package toolkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Do not include if empty with omitempty
}

// ReadJSON reads and decodes JSON data from an HTTP request body into the provided 'data' object.
// It ensures the JSON is properly formatted, validates its size, and handles various error scenarios.
func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// Check if the payload is of permitted size
	maxBytes := 1024 * 1024 // 1 Mg
	if t.MaxJSONSize != 0 {
		maxBytes = t.MaxJSONSize
	}

	// Read request of the body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decode the body
	decodedBody := json.NewDecoder(r.Body)

	// Check if we should process JSON with unknown fields
	if !t.AllowUnknownFields {
		decodedBody.DisallowUnknownFields()
	}

	// Decode data
	err := decodedBody.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			// If there's a syntax error in the JSON, report the position of the error
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			// If the body is incomplete, return a malformed JSON error
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			// If there's a type mismatch, report which field is problematic
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			// If there's no specific field, report the character offset where the type error occurred
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			// If the body is empty, return an error indicating that the body must not be empty
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			// If there is an unknown field in the JSON, return an error indicating which field is unknown
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			// If the body exceeds the allowed size, return an error with the size limit
			return fmt.Errorf("body must not be larger %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			// If unmarshalling fails for any reason, return the error message
			return fmt.Errorf("error unmarshalling JSON %s", err.Error())
		default:
			return err
		}
	}

	// Ensure that only one JSON file is received
	// Attempt to decode an empty struct after the initial JSON decoding
	// If more JSON data is found, return an error indicating multiple JSON objects
	err = decodedBody.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one JSON value")
	}

	return nil
}

// WriteJSON() writes a JSON response with provided status, data and an optional custom header
func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	// Attempt to marshal data into JSON
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Check if a custom header should be set
	if len(headers) > 0 {
		for indx, hdr := range headers[0] {
			w.Header()[indx] = hdr
		}
	}

	// Set Content-Type and provided status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON() takes in an error and an optional status code, and sends a JSON error message
func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	// Set a default status
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var JSONPayload JSONResponse
	JSONPayload.Error = true
	JSONPayload.Message = err.Error()

	return t.WriteJSON(w, statusCode, JSONPayload)
}
