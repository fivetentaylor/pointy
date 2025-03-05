package dynamo

import (
	"crypto/sha256"
	"encoding/base64"
	"reflect"
)

func MustListAllDynamodbavs(obj interface{}) []string {
	attrs := []string{}
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		attrs = append(attrs, field.Tag.Get("dynamodbav"))
	}
	return attrs
}

func hashAndEncode(input []byte, outputLength int) string {
	hash := sha256.Sum256(input)

	if outputLength > len(hash) {
		outputLength = len(hash)
	}
	truncated := hash[:outputLength]

	encoded := base64.RawURLEncoding.EncodeToString(truncated)

	return encoded
}
