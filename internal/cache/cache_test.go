package cache

import (
	"sync"
	"testing"
)

func TestConcurrentCache_RaceDetector(t *testing.T) {
	c, _ := New(100)
	var wg sync.WaitGroup

	// Fire up 1,000 concurrent workers smashing the cache
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			c.Write("key", val)
		}(i)

		go func() {
			defer wg.Done()
			c.Access("key")
		}()
	}
	wg.Wait()
}
