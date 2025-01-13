package toolkit

import "testing"

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
