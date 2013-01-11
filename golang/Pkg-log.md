Pkg log
----

Package log实现了简单的日志功能. 定义了类型Logger, 用于输出格式化日志.  
同时, 预定义了名为std的Logger对象, 用于输出日志到stderr.  
包中的Fatal方法将在输出日志消息后调用os.Exit(1), 而Panic方法将在输出日志消息后调用panic.

### Logger
Logger可在多个goroutines中并发使用, 其内部使用sync.Mutex进行同步.
```go
type Logger struct {
	// 同步锁
	mu     sync.Mutex
	// 日志消息前缀
	prefix string
	// 用于控制日期, 时间, 文件名等信息的输出
	flag   int
	// 输出的目标
	out    io.Writer
	buf    []byte
}
```
flag可取以下值:
```go
const (
	// 输出日期: 2009/01/23
	Ldate         = 1 << iota
	// 输出时间: 01:23:23
	Ltime
	// 输出毫秒值: 01:23:23.123123
	Lmicroseconds
	// 输出文件的绝对路径名和行数: /a/b/c/d.go:23
	Llongfile
	// 输出文件名和行数: d.go:23, 设置了此flag, 会覆盖Llongfile
	Lshortfile
	// 同时输出日期和时间
	LstdFlags     = Ldate | Ltime
)
```
若想同时输出日期和完整的文件名, 可将flag设置为: Ldate | Llongfile.

### Logger.New
创建Logger对象.
```go
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag}
}
```

### Logger.Output
输出日志消息.
```go
// calldepth指定调用深度, s表示要输出的日志消息
func (l *Logger) Output(calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string // 文件名
	var line int // 行数
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// runtime.Caller用于返回调用堆栈信息, 该操作是昂贵的, 因此先释放锁, 避免阻塞其他goroutine的日志输出
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		// 考虑获取堆栈信息失败的情况
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	// 重置缓冲区
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	// 将日志消息也写入缓冲区
	l.buf = append(l.buf, s...)
	// 如果最后一个字符不是换行符, 则在末尾加入'\n'
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	// 写入io.Writer
	_, err := l.out.Write(l.buf)
	return err
}
```
formatHeader方法将prefix, flag等信息写入缓冲区.
```go
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	// 写入prefix
	*buf = append(*buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {	// 写入日期
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 { // 写入时间
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {	// 写入文件名和行数
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}
```
综上, 每一个日志消息包含3个部分, 从头到尾依次是: 
+ prefix, 创建Logger时指定.
+ flag对应的信息, 创建Logger时指定. 根据flag的不同取值, 此处可能包含日期, 时间, 文件名等信息.
+ 真实的日志消息. 每次调用Output时指定.

### 格式化输出方法
这些方法都在内部调用Logger.Output, 其calldepth设定为2.
```go
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Print(v ...interface{}) { l.Output(2, fmt.Sprint(v...)) }

func (l *Logger) Println(v ...interface{}) { l.Output(2, fmt.Sprintln(v...)) }

func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
	// 结束程序
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}
```

### Logger.Flags和Logger.SetFlags
返回或设置Logger的flag.
```go
func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}
```

### Logger.Prefix和Logger.SetPrefix
返回或设置Logger的prefix.
```go
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}
```

### Logger.SetOutput
设置Logger的输出目标.
```go
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}
```

### std
std是package内部定义的一个Logger.
```go
var std = New(os.Stderr, "", LstdFlags)
```



links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
