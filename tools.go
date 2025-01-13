package toolkit

import (
	"crypto/rand" // cryptographically secure random number generator
)

const randomStrSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

// Tools is the type used to instantiate this module.
// Any variable of this type will have access to all the methods
// with the receiver *Tools.
type Tools struct{}

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
