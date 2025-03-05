package testutils

import (
	"encoding/json"
	"testing"
)

func JsonMarshalString(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
