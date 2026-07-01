package circuitbreaker

import (
	"sync"
	"sync/atomic"
	"time"
)

type bucket struct {
	total    int64
	failures int64
}

type slidingWindow struct {
	mu       sync.Mutex
	size     time.Duration
	interval time.Duration
	buckets  []*bucket
	head     int
	headTime time.Time
}

func newSlidingWindow(size time.Duration) *slidingWindow {
	numBuckets := 10
	interval := size / time.Duration(numBuckets)
	if interval <= 0 {
		interval = 1 * time.Second
	}

	buckets := make([]*bucket, numBuckets)
	for i := range buckets {
		buckets[i] = &bucket{}
	}

	return &slidingWindow{
		size:     size,
		interval: interval,
		buckets:  buckets,
		headTime: time.Now(),
	}
}

func (w *slidingWindow) currentBucket() *bucket {
	now := time.Now()
	w.mu.Lock()
	defer w.mu.Unlock()

	elapsed := now.Sub(w.headTime)
	steps := int(elapsed / w.interval)

	if steps >= len(w.buckets) {
		for i := range w.buckets {
			w.buckets[i] = &bucket{}
		}
		w.head = 0
		w.headTime = now
		return w.buckets[0]
	}

	for i := 0; i < steps; i++ {
		w.head = (w.head + 1) % len(w.buckets)
		w.headTime = w.headTime.Add(w.interval)
		w.buckets[w.head] = &bucket{}
	}

	return w.buckets[w.head]
}

func (w *slidingWindow) addSuccess() {
	b := w.currentBucket()
	atomic.AddInt64(&b.total, 1)
}

func (w *slidingWindow) addFailure() {
	b := w.currentBucket()
	atomic.AddInt64(&b.total, 1)
	atomic.AddInt64(&b.failures, 1)
}

func (w *slidingWindow) totalRequests() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	var total int64
	for _, b := range w.buckets {
		total += atomic.LoadInt64(&b.total)
	}
	return total
}

func (w *slidingWindow) successCount() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	var successes int64
	for _, b := range w.buckets {
		total := atomic.LoadInt64(&b.total)
		failures := atomic.LoadInt64(&b.failures)
		successes += total - failures
	}
	return successes
}

func (w *slidingWindow) failureRate() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	var total, failures int64
	for _, b := range w.buckets {
		total += atomic.LoadInt64(&b.total)
		failures += atomic.LoadInt64(&b.failures)
	}
	if total == 0 {
		return 0
	}
	return float64(failures) / float64(total)
}

func (w *slidingWindow) reset() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i := range w.buckets {
		w.buckets[i] = &bucket{}
	}
	w.headTime = time.Now()
}
