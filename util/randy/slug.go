package randy

import (
	"math/rand"
)

// Alphabets
var (
	AlphabetAlpha            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	AlphabetAlphaNumeric     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	AlphabetUppercaseNumeric = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	AlphabetLowercaseNumeric = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	AlphabetNumeric          = []rune("0123456789")
)

// Defaults
var (
	DefaultSlugLength = uint8(8)
	DefaultLength     = uint8(32)
	DefaultAlphabet   = AlphabetAlphaNumeric
)

// Slug creates a random string of given length.
func Slug(length uint8) string {
	if length == 0 {
		length = DefaultSlugLength
	}

	return String(nil, length)
}

// String creates a random string of given length using the given alphabet.
func String(alphabet []rune, length uint8) string {
	if len(alphabet) == 0 {
		alphabet = DefaultAlphabet
	}

	if length == 0 {
		length = DefaultLength
	}

	b := make([]rune, length)

	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}

func Numeric(length uint8) string {
	return String(AlphabetNumeric, length)
}
