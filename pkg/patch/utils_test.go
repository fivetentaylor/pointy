package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentenceBounds(t *testing.T) {
	i, j := SentenceBounds("", 0)
	assert.Equal(t, i, j, "i should equal j")

	i, j = SentenceBounds("Hello World!", 0)
	assert.Equal(t, i, 0, "i should equal 0")
	assert.Equal(t, j, 11, "j should equal 11")

	i, j = SentenceBounds("Hello. World! Cool\n", 0)
	assert.Equal(t, i, 0, "i should equal 0")
	assert.Equal(t, j, 5, "j should equal 5")

	i, j = SentenceBounds("Hello. World! Cool\n", 6)
	assert.Equal(t, i, 6, "i should equal 6")
	assert.Equal(t, j, 12, "j should equal 12")
}
