package utils

import (
	"math/rand"
	"strings"
	"time"
)

func RandomRune(rand *rand.Rand) rune {
	// Define a range of Unicode code points; this example uses a basic range
	// You might want to define a more specific range depending on your needs
	min := rune(0x10000) // Start of Supplementary Multilingual Plane
	max := rune(0x1FFFF) // End of the first Supplementary Plane
	return rune(rand.Intn(int(max-min+1)) + int(min))
}

func RandomString(rand *rand.Rand, length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteRune(RandomRune(rand))
	}
	return result.String()
}

var lettersAndNumbers = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var lettersAndNumbersLength = len(lettersAndNumbers)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generates a random string of length n with characters from a-Z and 0-9
func RandomSafeString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(lettersAndNumbersLength)]
	}
	return string(b)
}
