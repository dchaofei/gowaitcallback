waitcallback-go
==========

应用场景：有一些需要异步返回结果的请求却需要同步去等待

文档
=============

完整的文档 [Godoc](https://godoc.org/github.com/dchaofei/waitcallback-go)

示例
=======

```go
waitCallback := NewCallback()
wg := sync.WaitGroup{}
wg.Add(2)
key := "ding"
value := "dong"
go func () {
	defer wg.Done()
	got := waitCallback.PushWait(key)
	fmt.Println(got)
}()
time.Sleep(time.Millisecond)
go func () {
	defer wg.Done()
	waitCallback.Resolve(key, value)
}()
wg.Wait()
```