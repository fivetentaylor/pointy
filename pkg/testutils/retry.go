package testutils

import (
	"time"
)

func RetryOnError(retries int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < retries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(sleep)
	}
	return err
}

func RetryUntilTrue(retries int, sleep time.Duration, fn func() bool) bool {
	var b bool
	for i := 0; i < retries; i++ {
		if b = fn(); b == true {
			return true
		}
		time.Sleep(sleep)
	}
	return false
}
