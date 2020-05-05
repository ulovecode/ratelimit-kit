package fixed_window_counter

import (
	"sync"
	"sync/atomic"
	"time"

	ratelimit_kit "github.com/ulovecode/ratelimit-kit"
)

var (
	once sync.Once
)

var _ ratelimit_kit.RateLimiter = &fixedWindowCounter{}

type fixedWindowCounter struct {
	snippet         time.Duration
	currentRequests int32
	allowRequests   int32
}

func New(snippet time.Duration, allowRequests int32) *fixedWindowCounter {
	return &fixedWindowCounter{snippet: snippet, allowRequests: allowRequests}
}

func (l *fixedWindowCounter) Take() error {
	once.Do(func() {
		go func() {
			for {
				select {
				case <-time.After(l.snippet):
					atomic.StoreInt32(&l.currentRequests, 0)
				}
			}
		}()
	})

	curRequest := atomic.LoadInt32(&l.currentRequests)
	if curRequest >= l.allowRequests {
		return ratelimit_kit.ErrExceededLimit
	}
	if !atomic.CompareAndSwapInt32(&l.currentRequests, curRequest, curRequest+1) {
		return ratelimit_kit.ErrExceededLimit
	}
	return nil
}
