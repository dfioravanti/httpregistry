package httpregistry_test

import (
	"encoding/json"
	"fmt"
)

// mustMarshalJSON tries to marshal v into JSON and panics if it cannot
func mustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("body cannot be marshaled to JSON: %s", err))
	}
	return b
}
