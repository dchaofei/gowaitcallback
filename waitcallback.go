// 应用场景：有一些需要异步返回结果的请求却需要同步去等待
package waitcallback

import (
	"context"
	"sync"
)

// WaitCallback 等待Callback的结构体
type WaitCallback struct {
	valMap sync.Map
}

type wrapValue struct {
	ch    chan struct{}
	value interface{}
}

func NewCallback() *WaitCallback {
	return &WaitCallback{}
}

func (c *WaitCallback) pushWait(ctx context.Context, key interface{}) interface{} {
	var wrValue *wrapValue
	value, ok := c.valMap.Load(key)
	if ok {
		wrValue = value.(*wrapValue)
	} else {
		wrValue = &wrapValue{
			ch: make(chan struct{}),
		}
		c.valMap.Store(key, wrValue)
	}
	select {
	case <-wrValue.ch:
		return wrValue.value
	case <-ctx.Done():
		return nil
	}
}

// PushWait 阻塞一直到key被回调返回值
func (c *WaitCallback) PushWait(key interface{}) interface{} {
	return c.pushWait(context.Background(), key)
}

// PushWaitWithContext 阻塞到key被回调返回值或者context超时结束
func (c *WaitCallback) PushWaitWithContext(ctx context.Context, key interface{}) interface{} {
	return c.pushWait(ctx, key)
}

// Resolve 填充key对应的值并利用close(ch)广播通知
func (c *WaitCallback) Resolve(key interface{}, value interface{}) {
	v, ok := c.valMap.Load(key)
	if !ok {
		return
	}
	wrValue := v.(*wrapValue)
	wrValue.value = value
	c.valMap.Delete(key)
	close(wrValue.ch)
}
