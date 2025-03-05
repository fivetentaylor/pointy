package testutils

import "os"

// LiveEnabled returns true if the LIVE_TEST env var is set
// If you set LIVE_TEST=true, make sure you also set all the LLM keys + freeplay to real values
// or the tests will fail
func LiveEnabled() bool {
	return os.Getenv("LIVE_TEST") == "true"
}
