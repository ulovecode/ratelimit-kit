package token_bucket

import (
	"sync"
	"time"

	ratelimit_kit "github.com/ulovecode/ratelimit-kit"
)

var (
	once sync.Once
)
var _ ratelimit_kit.RateLimiter = &tokenBucket{}

type tokenBucket struct {
	snippet       time.Duration
	token         chan struct{}
	allowRequests int32
}

func New(snippet time.Duration, allowRequests int32) *tokenBucket {
	return &tokenBucket{snippet: snippet, token: make(chan struct{}, allowRequests), allowRequests: allowRequests}
}

func (t *tokenBucket) Take() error {
	once.Do(func() {
		go func() {
			for {
				select {
				case <-time.After(time.Duration(t.snippet.Nanoseconds() / int64(t.allowRequests))):
					t.token <- struct{}{}
				}
			}
		}()
	})
	select {
	case <-t.token:
		return nil
	default:
	}
	return ratelimit_kit.ErrExceededLimit
}
