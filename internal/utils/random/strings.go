package random

import "crypto/rand"

const characters = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

// GenerateString generate a random string based on a given size
func GenerateString(size int) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}
