package eria

import (
	"sync"
	"testing"
)

// The remote Things connect in parallel: they must share one consumer.
func TestConsumerIsShared(t *testing.T) {
	var wg sync.WaitGroup
	got := make([]any, 8)
	for i := range got {
		wg.Add(1)
		go func() { defer wg.Done(); got[i] = Consumer() }()
	}
	wg.Wait()
	for _, c := range got {
		if c != got[0] {
			t.Fatal("two consumers made")
		}
	}
}
