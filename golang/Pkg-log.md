Pkg log
----

Package log实现了简单的日志功能. 定义了类型Logger, Logger用于输出格式化日志.  
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
	flag   int
	// 输出的目标
	out    io.Writer
	buf    []byte
}
```


links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
