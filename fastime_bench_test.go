package fastime

import (
	"context"
	"testing"
	"time"
)

// BenchmarkFastime
func BenchmarkFastime(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	b.Cleanup(cancel)

	t := New()
	t.StartTimerD(ctx, time.Millisecond*5)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			t.Now()
		}
	})
}

// BenchmarkTime
func BenchmarkTime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			time.Now()
		}
	})
}
