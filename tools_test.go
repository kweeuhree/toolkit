package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	tests := []struct {
		name   string
		length int
	}{
		{"Ten", 10},
		{"Twenty", 20},
		{"One", 1},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			s := testTools.RandomString(entry.length)

			if len(s) != entry.length {
				t.Errorf("expected length %d, received %d", entry.length, len(s))
			}
		})
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	tests := []struct {
		name          string
		allowedTypes  []string
		renameFile    bool
		errorExpected bool
	}{
		{"allowed no rename", []string{"image/jpeg", "image/png"}, false, false},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Set up a pipe to avoid buffering while simulating a multipart file upload
			pr, pw := io.Pipe()
			// Create a multipart writer
			mpWriter := multipart.NewWriter(pw)

			// execute a goroutine in the background
			// use a wait group to ensure that things occur in a particular sequence
			wg := sync.WaitGroup{}
			// Add one to a wait group
			wg.Add(1)

			// Create a go func to run concurrently with the program
			go func() {
				// Close the writer when the function is finished
				defer mpWriter.Close()
				// Decrement the wait group by one when the function is finished
				defer wg.Done()

				// Create the form data field 'file'
				part, err := mpWriter.CreateFormFile("file", "./testdata/img.png")
				if err != nil {
					t.Error(err)
				}

				file, err := os.Open("./testdata/img.png")
				if err != nil {
					t.Error(err)
				}

				// Close the file to avoid resource leaks
				defer file.Close()

				// Decode the image
				img, _, err := image.Decode(file)
				if err != nil {
					t.Error("error decoding image", err)
				}

				err = png.Encode(part, img)
			}()

			// Read from the pipe which receives data
			// Create a request with a pipe reader
			req := httptest.NewRequest(http.MethodPost, "/", pr)
			// Set the correct content type for whatever type the payload is
			req.Header.Add("Content-Type", mpWriter.FormDataContentType())

			var testTools Tools
			testTools.AllowedFileTypes = entry.allowedTypes
			// Call UploadFiles with the pipe reader request, save to 'uploads' folder
			uploadedFiles, err := testTools.UploadFiles(req, "./testdata/uploads/", entry.renameFile)
			// Fail the test if the error is not nil and was not expected
			if err != nil && !entry.errorExpected {
				t.Errorf("expected no error, but received %+v", err)
			}

			// Build a string with the new file name
			fileStr := fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)

			// Check if the file did not get uploaded
			if !entry.errorExpected {
				if _, err := os.Stat(fileStr); os.IsNotExist(err) {
					t.Errorf("expected file to exist: %s", err.Error())
				}

				// Clean up
				_ = os.Remove(fileStr)
			}

			// If error is expected and not received, log it
			if entry.errorExpected && err == nil {
				t.Errorf("expected an error, but none received")
			}

		})
	}
}
