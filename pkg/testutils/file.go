package testutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func File(t *testing.T, path string) (*os.File, int64) {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)
	path = filepath.Join(basePath, "../../", path)

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	return file, fileInfo.Size()
}
