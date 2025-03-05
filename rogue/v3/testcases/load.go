package testcases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	v3 "github.com/teamreviso/code/rogue/v3"
)

type Failable interface {
	Fatalf(string, ...interface{})
}

type LogFail struct{}

func (l LogFail) Fatalf(s string, args ...interface{}) {
	log.Fatalf(s, args...)
}

func Load(t Failable, name string) *v3.Rogue {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to retrieve caller information")
	}
	dir := filepath.Dir(filename)

	// Change the base directory to "../../testcases" relative to the current file directory
	baseDir := filepath.Join(dir, "../../testcases")

	// Construct the path with the new base directory
	path := filepath.Join(baseDir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}

	var doc v3.Rogue
	err = json.Unmarshal(data, &doc)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	return &doc
}

type TestCase struct {
	Name string
	Doc  *v3.Rogue
}

func LoadAll(t Failable) []TestCase {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to retrieve caller information")
	}
	dir := filepath.Dir(filename)

	// Change the base directory to "../../testcases" relative to the current file directory
	baseDir := filepath.Join(dir, "../../testcases")

	files, err := os.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %s", err)
	}

	out := make([]TestCase, 0, len(files))
	for _, dirEntry := range files {
		name := dirEntry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		path := filepath.Join(baseDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("[%s] Failed to read file: %s", name, err)
		}

		var sRogue v3.SerializedRogue
		if err := json.Unmarshal(data, &sRogue); err != nil {
			t.Fatalf("[%s] Failed to unmarshal JSON: %s", name, err)
		}

		r := v3.NewRogueForQuill("auth0")
		snapshot := v3.SnapshotOp{Snapshot: &sRogue}

		_, err = r.MergeOp(snapshot)
		if err != nil {
			t.Fatalf("[%s] Failed to merge snapshot: %s", name, err)
		}

		out = append(out, TestCase{
			Name: name,
			Doc:  r,
		})
	}

	return out
}

func ConvertAll() {
	t := LogFail{}
	cases := LoadAll(t)
	for _, c := range cases {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatalf("Failed to retrieve caller information")
		}
		dir := filepath.Dir(filename)

		baseDir := filepath.Join(dir, "../../testcases")
		fullPath := filepath.Join(baseDir, c.Name)
		fmt.Printf("overwrite: %s\n", fullPath)

		data, err := json.Marshal(c.Doc)
		if err != nil {
			t.Fatalf("[%s] Failed to marshal JSON: %s", c.Name, err)
		}

		file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("[%s] Failed to open file: %s", c.Name, err)
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			t.Fatalf("[%s] Failed to write file: %s", c.Name, err)
		}
	}
}
