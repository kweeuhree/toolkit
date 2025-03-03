package toolkit

import (
	"crypto/rand" // cryptographically secure random number generator
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// Tools is the type used to instantiate this module.
// Any variable of this type will have access to all the methods
// with the receiver *Tools.
type Tools struct {
	MaxFileSize        int      // Specify the max size of a file permitted for uploading
	AllowedFileTypes   []string // Specify the file types to be permitted for uploading
	MaxJSONSize        int      // Specify the max size of a JSON payload
	AllowUnknownFields bool     // Permit the unknown fields
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

// CreateNewDirectory() creates a new directory if it does not exist
func (t *Tools) CreateNewDirectory(path string) error {
	// Set a mode to set a regular directory with following permissions:
	// - Owner: read, write, execute
	// - Group: read, execute
	// - Others: read, execute
	const mode = 0755
	// Check if the directory already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err = os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

// Slugify() takes in a string and replaces all but letters and numbers with hyphens
func (t *Tools) Slugify(str string) (string, error) {
	trimmmed := strings.Trim(str, " ")
	if len(trimmmed) == 0 {
		return "", errors.New("empty string not permitted")
	}

	// Define a pattern for a string of any length containing any letter or digit
	regex := regexp.MustCompile(`[^a-z\d]+`)

	slug := regex.ReplaceAllString(strings.ToLower(trimmmed), "-")

	// Ensure the slug is not empty and does not end or start with a dash
	slug = strings.Trim(slug, "-")
	if len(slug) == 0 {
		return "", errors.New("resultant string is empty")
	}

	return slug, nil
}

// DownloadStaticFile() downloads a file from the server to the local users machine
func (t *Tools) DownloadStaticFile(w http.ResponseWriter, r *http.Request, dirPath, fileName, displayName string) {
	// Construct the file path by joining the provided directory path and file name
	filePath := path.Join(dirPath, fileName)

	// Set the response header to indicate a file attachment with the specified display name
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))

	// Serve the file to the user, prompting a download
	http.ServeFile(w, r, filePath)
}
