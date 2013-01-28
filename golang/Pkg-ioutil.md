Pkg ioutil
----

ioutil包提供了一些io操作的便利方法.

### ReadAll
从指定的io.Reader对象中读取数据, 直到发生error或者到达文件末尾, 返回读取的所有数据.
```go
func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}

// 第二个参数表示初始申请的内存块的大小
// 如果该值过大会导致内存的浪费; 如果过小, 则会导致频繁的重新申请内存
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	defer func() {
		// readAll方法试图将io.Reader(通常是文件)的所有内容读进内存
		// 如果待读取的数据量过大, 需要申请的内存可能超过操作系统限制
		// 此时将导致buf.ReadFrom方法发生panic, recover该panic时能得到bytes.ErrTooLarge
		// 在defer调用中, 只处理上述提到的panic, 将其转换为err返回给调用方, 而保持其他未知的panic
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

### ReadFile
ReadFile和ReadAll方法类似, 区别在于ReadFile读取的是指定文件的所有内容.
```go
func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	// 接下来试图根据文件的大小, 给定一个合适的capacity参数(在调用readAll方法时会用到). 
	// 虽然FileInfo返回的文件大小不一定准确, 但仍然值得一试
	var n int64

	if fi, err := f.Stat(); err == nil {
		// Don't preallocate a huge buffer, just in case.
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	return readAll(f, n+bytes.MinRead)
}
```

### Readdir
获取指定目录下的文件列表, 并根据文件名进行排序.
```go
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	// 排序
	sort.Sort(byName(list))
	return list, nil
}

// byName实现了sort.Interface接口所需的Len, Less, Swap方法.
// 以方便根据文件名进行排序.
type byName []os.FileInfo
func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
```

### NopCloser
将io.Reader对象包装成io.ReadCloser对象.
```go
func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{r}
}

// nopCloser匿名内嵌io.Reader, 因此nopCloser对象实现了io.Reader接口
// 又因为nopCloser对象具有Close方法, 因此nopCloser对象实现了io.ReadCloser接口
type nopCloser struct {
	io.Reader
}
// Close方法什么也不做, 只是为了实现接口而存在
func (nopCloser) Close() error { return nil }
```

### Discard
Discard是一个io.Writer对象, 该对象相当于bash中的/dev/null, 所有写往Discard的数据都将被丢弃, 并保证写入操作肯定成功.
```go
var Discard io.Writer = devNull(0)
type devNull int
// 实现io.Writer接口. all Write calls succeed without doing anything.
func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}
```

### TempFile
使用指定的文件名前缀, 在指定的目录中, 创建临时文件.
```go
func TempFile(dir, prefix string) (f *os.File, err error) {
	// 如果dir为空, 使用默认的临时目录, 在linux下, 该目录为/tmp
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix())
		// 同时指定了os.O_CREATE, os.O_EXCL调用OpenFile时, 如果文件已存在, 该方法将返回error
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			// 冲突次数超过10时, 调整随机种子
			if nconflict++; nconflict > 10 {
				rand = reseed()
			}
			continue
		}
		break
	}
	return
}
```


















links
-----
+ [目录](../golang)
+ 上一节: [Pkg bufio](Pkg-bufio.md)
+ 下一节: [Pkg list](Pkg-list.md)
