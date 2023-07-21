# 测试 chan source-filter-sink

如果想测试 `zero-pump`，可以使用以下代码来替换 `Run`：

```go
func (s preparedServer) Run() error {
	source := ext.NewChanSource(tickerChan(time.Second * 1))

	filter := flow.NewMap(addUTC, 1)

	sink := ext.NewStdoutSink()

	source.Via(filter).To(sink)
	return nil
}

func tickerChan(repeat time.Duration) chan interface{} {
	ticker := time.NewTicker(repeat)
	oc := ticker.C
	nc := make(chan interface{})
	go func() {
		for range oc {
			nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
		}
	}()
	return nc
}
```
