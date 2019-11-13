package waitcallback

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWaitCallback_PushWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	waitCallback := NewCallback()
	const key = "ping"
	const value = "pong"
	go func() {
		defer wg.Done()
		got := waitCallback.PushWait(context.Background(), key)
		if got.(string) != value {
			t.Errorf("want %s got %s", value, got)
		}
	}()
	time.Sleep(time.Microsecond)
	go func() {
		defer wg.Done()
		waitCallback.Resolve(key, value)
	}()
	wg.Wait()
}

func TestWaitCallback_PushWaitWithTimeout(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	waitCallback := NewCallback()
	const key = "ping"
	go func() {
		defer wg.Done()
		timeoutCtx, _ := context.WithTimeout(context.Background(), time.Microsecond)
		got := waitCallback.PushWait(timeoutCtx, key)
		if got != nil {
			t.Errorf("want %s got %s", "nil", got)
		}
	}()
	wg.Wait()
}
