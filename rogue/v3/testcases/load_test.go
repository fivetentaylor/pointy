package testcases

import "testing"

func BenchmarkLoadAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadAll(b)
	}
}
