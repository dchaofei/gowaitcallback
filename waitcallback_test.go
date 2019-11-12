package waitcallback

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWaitCallback_PushWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	var wgPush sync.WaitGroup
	wgPush.Add(9)
	waitCallback := NewCallback()
	const key = "ping"
	const value = "pong"
	for i := 0; i < 9; i++ {
		go func() {
			wgPush.Done()
			defer wg.Done()
			got := waitCallback.PushWait(key)
			if got.(string) != value {
				t.Errorf("want %s got %s", value, got)
			}
		}()
	}
	wgPush.Wait()
	go func() {
		defer wg.Done()
		waitCallback.Resolve(key, value)
	}()
	wg.Wait()
}

func TestWaitCallback_PushWaitWithContext(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	waitCallback := NewCallback()
	const key = "ping"
	go func() {
		defer wg.Done()
		timeoutCtx, _ := context.WithTimeout(context.Background(), time.Microsecond)
		got := waitCallback.PushWaitWithContext(timeoutCtx, key)
		if got != nil {
			t.Errorf("want %s got %s", "nil", got)
		}
	}()
	wg.Wait()
}
