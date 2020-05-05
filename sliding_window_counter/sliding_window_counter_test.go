package sliding_window_counter

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
	tests := []struct {
		name string
		args args
		want *slidingWindowCounter
	}{
		{name: "0", args: struct {
			accuracy      time.Duration
			snippet       time.Duration
			allowRequests int32
		}{accuracy: time.Microsecond, snippet: time.Second, allowRequests: 100}, want: nil},
	}
	count := 0
	var f *slidingWindowCounter
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f = New(tt.args.accuracy, tt.args.snippet, tt.args.allowRequests)
			for range [10]struct{}{} {
				<-time.After(f.snippet)
				go func() {
					select {
					case <-time.After(f.accuracy):
						for range [50]struct{}{} {
							if err := f.Take(); err != nil {
								//t.Errorf("Take() error = %v", err)
							} else {
								//t.Log(i, "take once request")
								count++
							}
						}
					}
				}()
			}
		})
	}
	time.Sleep(2 * time.Second)
	if count != 100*5 {
		t.Fatal("count:", count)
	}
	t.Log("count:", count)
}
