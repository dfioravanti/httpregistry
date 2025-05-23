package httpregistry_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// mustMarshalJSON tries to marshal v into JSON and panics if it cannot
func mustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("body cannot be marshaled to JSON: %s", err))
	}
	return b
}

// generateRandomString generates a random string of a given length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
