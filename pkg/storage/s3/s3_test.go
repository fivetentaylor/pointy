package s3

import (
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file")
	}

	s3, err := NewS3()
	if err != nil {
		log.Fatalf("failed to create db: %s", err)
	}
	err = s3.CreateTestBucket()
	if err != nil {
		log.Infof("failed to create bucket: %s", err)
	}

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

func TestPutGet(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	s3, err := NewS3()
	require.NoError(t, err)

	objBytes := []byte("hello world, I am an s3 object")
	err = s3.PutObject(s3.Bucket, "test/path", "text/plain", objBytes)
	require.NoError(t, err)

	b, err := s3.GetObject(s3.Bucket, "test/path")
	require.NoError(t, err)

	ok, err := s3.Exists(s3.Bucket, "test/path")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = s3.Exists(s3.Bucket, "test/fake/path")
	require.NoError(t, err)
	assert.False(t, ok)

	require.Equal(t, b, objBytes)
	log.Info(string(b))
}
