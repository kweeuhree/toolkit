package toolkit

import (
	"crypto/rand" // cryptographically secure random number generator
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomStrSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

// Tools is the type used to instantiate this module.
// Any variable of this type will have access to all the methods
// with the receiver *Tools.
type Tools struct {
	MaxFileSize      int      // Specify the max size of a file permitted for uploading
	AllowedFileTypes []string // Specify the file types to be permitted for uploading
}

// RandomString() takes in an integer that defines length of random string.
// It uses randomStrSource as the source for the return string.
// Returns a string with the provided length.
func (t *Tools) RandomString(n int) string {
	// Preallocate a slice of runes with size 'n' to store the random characters.
	// Convert the random string source to a slice of runes for indexing.
	str, r := make([]rune, n), []rune(randomStrSource)
	// Loop over each index in the 'str' slice to fill it with a random character.
	for i := range str {
		// Generate a random prime number using a cryptographically secure random number generator.
		// The prime number's bit length is set to the length of the random string source (in bits).
		p, _ := rand.Prime(rand.Reader, len(r))
		// Convert the prime number to a uint64.
		x, y := p.Uint64(), uint64(len(r))

		// Select a character from the randomStrSource by taking the modulo of the prime number.
		// This ensures the index is within bounds of the source slice.
		str[i] = r[x%y]
	}

	// Convert the slice of runes into a string and return it.
	return string(str)
}

// UploadedFile is a struct used to save information about an uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile is a convenience method that calls UploadFiles
// Expectes only one file to be uploaded
func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// UploadFiles uploads one or more file to a specified directory, and gives the files a random name
// Returns a slice of with the newly named files, the original file names, file sizes, and
// a potential error. If the optional last parameter is set to true, the files will not be renamed
func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	// Rename by default
	renameFile := true

	// If the length of the variadic parameter is not zero, update renameFile
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	// Preallocate a slice to store the files
	var uploadedFiles []*UploadedFile

	// Assign MaxFileSize if it is not set
	if t.MaxFileSize == 0 {
		// Set a default limit
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	// Check for an error when parsing the request
	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("the uploaded file is too big")
	}

	// Check if any files are stored in the request
	for _, headers := range r.MultipartForm.File {
		for _, hdr := range headers {
			// Wrap defer in a function
			uploadedFiles, err = func(UploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				// Open the header
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				// Close in order to avoid resource leak
				defer infile.Close()

				// We need to look at the first 512 bytes to find out the type of file
				buff := make([]byte, 512)
				_, err = infile.Read(buff) // Read the bytes
				if err != nil {
					return nil, err
				}

				// Check to see if the file type is permitted
				// Assume that the file type is not allowed
				allowed := false
				fileType := http.DetectContentType(buff) // Get file type of the bytes

				// Check if the AllowedFileTypes was populated
				if len(t.AllowedFileTypes) > 0 {
					for _, f := range t.AllowedFileTypes {
						// If current file type equals one of the permitted file types...
						if strings.EqualFold(fileType, f) {
							// ...allow the file
							allowed = true
						}
					}
					// if AllowedFileTypes was not populated...
				} else {
					// ...allow all files
					allowed = true
				}

				// If allowed is still false, return an error
				if !allowed {
					return nil, errors.New("the uploaded file type is not permitted")
				}

				// Since we read the beginning of the file,
				// We have to go back to the beginning of the file
				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				// If its going to be renamed - generate a new name with original extension
				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				uploadedFile.OriginalFileName = hdr.Filename

				// Save to disk
				var outfile *os.File  // file we will write to
				defer outfile.Close() // close the file when the function exists

				// Write the file to the provided directory
				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}

					uploadedFile.FileSize = fileSize
				}

				// Append the file to the slice of uploadedFiles
				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil

				// give the function access to uploadedFiles
			}(uploadedFiles)

			// In case of error, return what was successfully uploaded
			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}
