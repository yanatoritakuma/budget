package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generates a random string of a given length.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a less random string in case of error
		return "couldnotgeneraterandomstring"
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}