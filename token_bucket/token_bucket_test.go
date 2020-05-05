package token_bucket

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		accuracy      time.Duration
		snippet       time.Duration
		allowRequests int32
	}
	var tests [1]struct {
		name string
		args args
		want *tokenBucket
	}
	count := 0
	var f *tokenBucket
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f = New(1*time.Second, 1000)
			f.Take()
			time.Sleep(1 * time.Second)
			for range [10]struct{}{} {
				go func() {
					<-time.After(f.snippet)
					for range [2000]struct{}{} {
						if err := f.Take(); err != nil {
							//t.Errorf("Take() error = %v", err)
						} else {
							//t.Log(i, "take once request")
							count++
						}
					}
				}()
			}
		})
		time.Sleep(2 * time.Second)
		if count != 1000 {
			t.Fatal("count:", count)
		}
		t.Log("count:", count)
	}

}
