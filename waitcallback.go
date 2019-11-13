// 应用场景：有一些需要异步返回结果的请求却需要同步去等待
package waitcallback

import (
	"context"
	"sync"
)

// WaitCallback 等待Callback的结构体
type WaitCallback struct {
	valMap map[interface{}]interface{}
	rwMu   sync.RWMutex
}

type wrapValue struct {
	ch    chan struct{}
	value interface{}
}

func NewCallback() *WaitCallback {
	return &WaitCallback{valMap: map[interface{}]interface{}{}}
}

// PushWait 阻塞到key被回调返回值或者context超时结束
func (c *WaitCallback) PushWait(ctx context.Context, key interface{}) interface{} {
	var wrValue *wrapValue
	c.rwMu.RLock()
	value, ok := c.valMap[key]
	c.rwMu.RUnlock()
	if ok {
		wrValue = value.(*wrapValue)
	} else {
		c.rwMu.Lock()
		// 双重检查
		value, ok := c.valMap[key]
		if ok {
			wrValue = value.(*wrapValue)
			c.rwMu.Unlock()
		} else {
			wrValue = &wrapValue{
				ch: make(chan struct{}),
			}
			c.valMap[key] = wrValue
			c.rwMu.Unlock()
		}
	}
	select {
	case <-wrValue.ch:
		return wrValue.value
	case <-ctx.Done():
		return nil
	}
}

// Resolve 填充key对应的值并利用close(ch)广播通知
func (c *WaitCallback) Resolve(key interface{}, value interface{}) {
	c.rwMu.RLock()
	v, ok := c.valMap[key]
	c.rwMu.RUnlock()
	if !ok {
		return
	}
	wrValue := v.(*wrapValue)
	wrValue.value = value
	delete(c.valMap, key)
	close(wrValue.ch)
}
