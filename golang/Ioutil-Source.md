Ioutil Source
----

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
```go
func (b *Reader) Read(p []byte) (n int, err error)
```
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
```go
func (b *Reader) ReadSlice(delim byte) (line []byte, err error)
```
读取并返回byte数据, 直到读取到指定的delim, 或者发生error, 或者缓冲区已满.  
如果缓冲区中的数据满足条件, 则直接返回:
```go
if i := bytes.IndexByte(b.buf[b.r:b.w], delim); i >= 0 {
	line1 := b.buf[b.r : b.r+i+1]
	b.r += i + 1
	return line1, nil
}
```
否则, 将数据读取到缓冲区后再进行判断:
```go
for {
	if b.err != nil {
		line := b.buf[b.r:b.w]
		b.r = b.w
		return line, b.readErr()
	}

	// 记录当前缓冲区中的字节数
	n := b.Buffered()
	b.fill()

	// 在新读取的数据中搜索
	if i := bytes.IndexByte(b.buf[n:b.w], delim); i >= 0 {
		line := b.buf[0 : n+i+1]
		b.r = n + i + 1
		return line, nil
	}

	// 判断缓冲区是否已满
	if b.Buffered() >= len(b.buf) {
		b.r = b.w
		return b.buf, ErrBufferFull
	}
}
```

### ReadLine
```go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
```
读取一行数据, 如果某行太长(超过缓冲区大小), 先返回开头部分, 并设置isPrefix为true, 其余部分将在随后的ReadLine中返回.  
先调用ReadSlice得到初步数据:
```go
line, err = b.ReadSlice('\n')
```
如果err是ErrBufferFull(缓冲区已满), 处理数据后返回给调用者:
```go
if err == ErrBufferFull {
	// 处理\r\n的情形
	if len(line) > 0 && line[len(line)-1] == '\r' {
		if b.r == 0 {
			// should be unreachable
			panic("bufio: tried to rewind past start of buffer")
		}
		// 将\r放回缓冲区, 并将其从line中删除, 以便下次调用ReadLine时一并处理'\r\n'
		b.r--
		line = line[:len(line)-1]
	}
	// 返回的isPrefix为true, 表明该行太长, line只是该行的一部分数据
	return line, true, nil
}
```
如果line的最后一个字节确实是'\n', 则将'\n'或者'\r\n'从line中删除后返回:
```go
if line[len(line)-1] == '\n' {
	drop := 1
	if len(line) > 1 && line[len(line)-2] == '\r' {
		drop = 2
	}
	line = line[:len(line)-drop]
}
```

### ReadBytes
```go
func (b *Reader) ReadBytes(delim byte) (line []byte, err error)
```
ReadBytes和ReadSlice的区别在于, ReadBytes不会因为缓冲区已满就停止搜索, 而是继续下去, 直到读取到指定的delim, 或者发生了error(通常是io.EOF).
```go
var frag []byte
var full [][]byte
err = nil

for {
	var e error
	frag, e = b.ReadSlice(delim)
	if e == nil { // 读取到指定的delim
		break
	}
	if e != ErrBufferFull { // 发生了ErrBufferFull之外的error
		err = e
		break
	}

	// 处理err为ErrBufferFull的情形
	// 将frag中的数据copy一份, 保存在full中.
	buf := make([]byte, len(frag))
	copy(buf, frag)
	full = append(full, buf)
}
```
如此以来, 数据就存储在full和frag中了, 接下来需要将2者中存储的数据合并到一个buffer中:
```go
// 计算buf所需的length
n := 0
for i := range full {
	n += len(full[i])
}
n += len(frag)

// 将数据从full和frag中copy到buf中, 注意copy的顺序
buf := make([]byte, n)
n = 0
for i := range full {
	n += copy(buf[n:], full[i])
}
copy(buf[n:], frag)
return buf, err
```

### ReadString
ReadString和ReadBytes类似, 区别在于ReadString返回的是string:
```go
func (b *Reader) ReadString(delim byte) (line string, err error) {
	bytes, e := b.ReadBytes(delim)
	return string(bytes), e
}
```
可以调用ReadString('\n')代替ReadLine, 此时不需要处理烦人的isPrefix标记, 因为对于ReadString('\n')来说, 不管行有多长, 都只会一次性返回.


Writer
----
bufio.Writer的定义如下:
```go
type Writer struct {
	err error
	// 缓冲区
	buf []byte
	// buf中有效数据个数
	n   int
	// 底层Writer
	wr  io.Writer
}
```

### 创建*Writer对象
可以通过以下2个function创建*bufio.Writer对象:
```go
// 指定缓冲区大小
func NewWriterSize(wr io.Writer, size int) *Writer
// 使用默认缓冲区大小, 相当于调用NewWriterSize(wr, defaultBufSize)
func NewWriter(wr io.Writer) *Writer
```

### Flush
将缓冲区中的数据写入相应的底层io.Writer:
```go
func (b *Writer) Flush() error {
	if b.err != nil {
		return b.err
	}
	// 缓冲区中没有数据
	if b.n == 0 {
		return nil
	}
	// 将缓冲区中的数据写入io.Writer
	n, e := b.wr.Write(b.buf[0:b.n])
	if n < b.n && e == nil {
		e = io.ErrShortWrite
	}
	if e != nil {
		// 写入了n个字节时, 将写入的部分从缓冲区中删除
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		// 设置相关标记
		b.n -= n
		b.err = e
		return e
	}
	// 清空缓冲区
	b.n = 0
	return nil
}
```

### Available
返回缓冲区中的空闲字节数:
```go
func (b *Writer) Available() int { return len(b.buf) - b.n }
```

### Buffered
返回缓冲区中的有效字节数:
```go
func (b *Writer) Buffered() int { return b.n }
```

### Write
*bufio.Writer具有Write方法, 因此*bufio.Writer实现了io.Writer接口.
```go
func (b *Writer) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// 需要写入大量数据, 且缓冲区中没有有效数据时, 直接将p写入底层io.Writer中
			n, b.err = b.wr.Write(p)
		} else {
			// 先将部分数据(b.Available()个字节)写入缓冲区中, 然后再刷新缓存
			n = copy(b.buf[b.n:], p)
			b.n += n
			// Flush时如果出错, 会将error存储在b.err中
			b.Flush()
		}
		nn += n
		// 截去已写入的部分
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	// 将剩余部分写入缓冲区, 此时len(p) <= b.Available(), 即缓冲区肯定能够容纳p中的数据
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}
```

### WriteByte
写入单个字节数据.
```go
func (b *Writer) WriteByte(c byte) error {
	if b.err != nil {
		return b.err
	}
	// 缓冲区中没有可用空间时, 尝试执行Flush操作, 如果Flush失败, 直接返回错误
	if b.Available() <= 0 && b.Flush() != nil {
		return b.err
	}
	// 将c写入缓冲区
	b.buf[b.n] = c
	b.n++
	return nil
}
```

### WriteString
写入string数据. 从源码可以看出, 相当于调用了`Write([]byte(s))`.








links
-----
+ [目录](../)
+ 上一节: [Bufio Source](Bufio-Source.md)
