package toolkit

import (
	"os"
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

func TestTools_CreateNewDirectory(t *testing.T) {
	tests := []struct {
		name    string
		dirPath string
		err     error
	}{
		{"Non-existing path", "./testFolder", nil},
		{"Non-existing path", "./testUploads", nil},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			var tools Tools
			// Call the CreateNewDirectory function
			err := tools.CreateNewDirectory(entry.dirPath)
			if err != nil && entry.err == nil {
				t.Errorf("expected no error, but received %+v", err)
			}

			// Remove the directory if it was created
			if _, err = os.Stat(entry.dirPath); err == nil {
				err = os.RemoveAll(entry.dirPath)
				if err != nil {
					t.Errorf("failed to remove path %s. %+v", entry.dirPath, err)
				}
			}
		})
	}
}

func TestTools_Slugify(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		err      bool
	}{
		{"Hello, World!", "hello-world", false},
		{"    !Hello && World!", "hello-world", false},
		// Cyrillic characters
		{"Привет, мир!Hello, World!", "hello-world", false},
		{"Привет, мир!", "", true},
		{"88GoLang!PyThon===Java?     TYPESCRIPT@   ", "88golang-python-java-typescript", false},
		{"        ", "", true},
		{"!!!!", "", true},
	}
	var tools Tools
	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result, err := tools.Slugify(entry.name)

			if result != entry.expected {
				t.Errorf("expected %s, but received %s", entry.expected, result)
			}

			if err != nil && !entry.err {
				t.Errorf("expected no error, but received %+v", err)
			}
		})
	}

}
