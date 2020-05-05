package leaky_barrel

import (
	"sync"
	"time"

	ratelimit_kit "github.com/ulovecode/ratelimit-kit"
)

var (
	once sync.Once
)

var _ ratelimit_kit.RateLimiter = &leakyBarrel{}

type leakyBarrel struct {
	snippet       time.Duration
	allowRequests int32
	barrelSize    chan struct{}
}

func New(snippet time.Duration, barrelSize int, allowRequests int32) *leakyBarrel {
	return &leakyBarrel{snippet: snippet, barrelSize: make(chan struct{}, int(allowRequests)/barrelSize), allowRequests: allowRequests}
}

func (t *leakyBarrel) Take() error {
	once.Do(func() {
		go func() {
			for {
				select {
				case <-time.After(time.Duration(t.snippet.Nanoseconds() / int64(t.allowRequests))):
					t.barrelSize <- struct{}{}
				}
			}
		}()
	})
	select {
	case <-t.barrelSize:
		return nil
	default:
	}
	return ratelimit_kit.ErrExceededLimit
}
