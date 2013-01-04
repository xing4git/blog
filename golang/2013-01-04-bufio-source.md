bufio包提供了带有缓冲功能的Reader和Writer, 应用了`装饰者模式`.  
bufio.Reader是对io.Reader的包装, 并且实现了io.Reader接口.  
类似的, bufio.Writer是对io.Writer的包装, 并且实现了io.Writer接口.


Reader
------
bufio.Reader的定义如下:
```go
type Reader struct {
	// 缓冲区
	buf          []byte
	// 底层Reader
	rd           io.Reader
	// buffer中有效数据的边界
	r, w         int
	// read过程中出现的error
	err          error
	lastByte     int
	lastRuneSize int
}
```
lastByte和lastRuneSize字段是为了支持unread操作而存在的.  
lastByte表示读取的最后一个字节数据, -1表示不可用状态. 注意lastByte的类型为int, 如果lastByte的类型为byte, -1是合法的字节数据值, 无法使用-1表示不可用状态. 将lastByte的类型设定为int则可以规避这个问题, byte强转为int之后, 取值范围是[0-255], 此时-1不再是合法值, 因此可以用来表示不可用状态.

### 创建*Reader对象
可以通过以下2个function创建*bufio.Reader对象:
```go
// 指定缓冲区大小
func NewReaderSize(rd io.Reader, size int) *Reader
// 使用默认缓冲区大小, 相当于调用NewReaderSize(rd, defaultBufSize)
func NewReader(rd io.Reader) *Reader
```
有趣的是, 在NewReaderSize中, 会判断rd是否已经是一个*bufio.Reader对象了:
```go
b, ok := rd.(*Reader)
if ok && len(b.buf) >= size {
	return b
}
```
如果rd已经是一个*bufio.Reader对象了, 且缓冲区大小也满足条件, 则直接返回rd自身, 以防止不必要的包装.

### Read
*bufio.Reader实现了io.Reader接口, 就是因为*bufio.Reader具有Read方法.  
如果缓冲区中没有数据:
```go
if b.w == b.r {
	if b.err != nil {
		return 0, b.readErr()
	}
	if len(p) >= len(b.buf) {
		// 读取数据量超过缓冲区大小时, 直接从底层的io.Reader对象读取并返回数据, 避免[]byte的copy
		// 注意, 这部分数据是没有进入缓冲区的
		n, b.err = b.rd.Read(p)
		if n > 0 {
			b.lastByte = int(p[n-1])
			b.lastRuneSize = -1
		}
		return n, b.readErr()
	}
	// 试图填满缓冲区
	b.fill()
	// 此时如果缓冲区仍然没有数据, 则说明肯定出现了error
	if b.w == b.r {
		return 0, b.readErr()
	}
}
```
接下来, 返回缓冲区中的数据:
```go
// 最多从缓冲区中取出n个字节
if n > b.w-b.r {
	n = b.w - b.r
}
copy(p[0:n], b.buf[b.r:])
// 改变标记字段
b.r += n
b.lastByte = int(b.buf[b.r-1])
b.lastRuneSize = -1
return n, nil
```
从代码可以看出, Read最多调用一次底层Reader的Read方法, 因此返回的n是有可能小于len(p)的.

### Buffered
返回当前缓冲区中的有效字节数.
```go
func (b *Reader) Buffered() int { return b.w - b.r }
```

### ReadSlice










