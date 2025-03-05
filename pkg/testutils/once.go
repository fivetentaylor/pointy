package testutils

import "sync"

type OnceRunner struct {
	once sync.Once
	wg   sync.WaitGroup
	mu   sync.Mutex
}
