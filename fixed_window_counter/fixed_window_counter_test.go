package fixed_window_counter

import (
	"testing"
	"time"
)

func Test_fixedWindowCounter_Take(t *testing.T) {
	tests := []struct {
		name    string
		fields  fixedWindowCounter
		wantErr error
	}{
		{name: "0",
			fields: fixedWindowCounter{
				snippet:         1 * time.Second,
				currentRequests: 0,
				allowRequests:   100,
			},
			wantErr: nil},
	}
	for _, tt := range tests {
		var f *fixedWindowCounter
		count := 0
		t.Run(tt.name, func(t *testing.T) {
			f = New(tt.fields.snippet, tt.fields.allowRequests)
			for range [10]struct{}{} {
				<-time.After(f.snippet)
				for range [1000]struct{}{} {
					go func() {
						for range [50]struct{}{} {
							if err := f.Take(); err != nil {
								//t.Errorf("Take() error = %v", err)
							} else {
								//t.Log(i, "take once request")
								count++
							}
						}
					}()

				}
			}
		})
		time.Sleep(2 * time.Second)
		if count != 100*10 {
			t.Fatal("count:", count)
		}
		t.Log("count:", count)
	}

}
