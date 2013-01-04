Package ioutil
----

ioutil包提供了一些io操作的便利方法.

### ReadAll
从指定的io.Reader对象中读取数据, 直到发生error或者到达文件末尾, 返回读取的所有数据.
```go
func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	defer func() {
		// (*bytes.Buffer).ReadFrom方法从io.Reader对象中读取数据到其缓冲区中
		// 如果bytes.Buffer内部的缓冲区不足以容纳所有数据, 则自动新建更大的缓冲区
		// 假如待读取的数据量实在过大, 
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
```



links
-----
+ [目录](../golang)
+ 上一节: [Package bufio](Package-bufio.md)
