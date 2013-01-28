Pkg expvar
----

Package expvar是用于获取public variables的标准接口. 通过http访问/debug/vars, 即可得到json格式的public variables.  
expvar默认publish了2个variable: cmdline和memstats, 分别用于获取程序启动时的命令行参数和当前内存使用情况.  
可根据自己的需要publish其他variable.

### init
初始化操作.
```go
func init() {
	http.HandleFunc("/debug/varss", expvarHandler)
	// publish cmdline and memstats
	Publish("cmdline", Func(cmdline))
	Publish("memstats", Func(memstats))
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	// 输出所有已publish的variable
	Do(func(kv KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

// 存储已publish的variables
var (
	mutex sync.RWMutex
	// 每个publish的variable需要实现Var接口, Var接口只定义了一个方法: String
	vars  map[string]Var = make(map[string]Var)
)

func Do(f func(KeyValue)) {
	mutex.RLock()
	defer mutex.RUnlock()
	for k, v := range vars {
		f(KeyValue{k, v})
	}
}
```

### Puslish
将variable添加到`expvar.vars`中.
```go
func Publish(name string, v Var) {
	mutex.Lock()
	defer mutex.Unlock()
	// name已存在时, 会引起panic
	if _, existing := vars[name]; existing {
		log.Panicln("Reuse of exported var name:", name)
	}
	vars[name] = v
}
```

### Var
所有需要publish的variable都必须实现Var接口.
```go
type Var interface {
	String() string
}
```
package expvar提供了一些实现Var接口的类型.

### Int, Float, String
Int是对int的包装, 实现Var接口, 并提供了一些原子操作.
```go
type Int struct {
	i  int64
	mu sync.RWMutex
}
// 实现Var接口
func (v *Int) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return strconv.FormatInt(v.i, 10)
}
// 原子操作
func (v *Int) Add(delta int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.i += delta
}
// 原子操作
func (v *Int) Set(value int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.i = value
}
```
类似的, Float是对float64的包装, 实现Var接口, 并提供了一些原子操作.
```go
type Float struct {
	f  float64
	mu sync.RWMutex
}
func (v *Float) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return strconv.FormatFloat(v.f, 'g', -1, 64)
}
func (v *Float) Add(delta float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.f += delta
}
func (v *Float) Set(value float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.f = value
}
```
String也是如此.
```go
type String struct {
	s  string
	mu sync.RWMutex
}
func (v *String) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return strconv.Quote(v.s)
}
func (v *String) Set(value string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.s = value
}
```

### Map
Map是对`map[string]Var`的包装, 同时其自身也实现了Var接口.
```go
type Map struct {
	m  map[string]Var
	mu sync.RWMutex
}

// map中的一个entry
type KeyValue struct {
	Key   string
	Value Var
}

// 实现了Var接口. 
func (v *Map) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	var b bytes.Buffer
	fmt.Fprintf(&b, "{")
	first := true
	for key, val := range v.m {
		if !first {
			fmt.Fprintf(&b, ", ")
		}
		fmt.Fprintf(&b, "\"%s\": %v", key, val)
		first = false
	}
	fmt.Fprintf(&b, "}")
	return b.String()
}

func (v *Map) Init() *Map {
	v.m = make(map[string]Var)
	return v
}

// ... 省略了一些方法

// 对Map中每对KeyValue调用f
func (v *Map) Do(f func(KeyValue)) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for k, v := range v.m {
		f(KeyValue{k, v})
	}
}
```

### Func
以上几种Var类型具有一个共同点: 他们都是`静态值`.  
Func可用于动态的获取某个值.
```go
type Func func() interface{}
// 为Func实现Var接口
func (f Func) String() string {
	v, _ := json.Marshal(f())
	return string(v)
}
```
Publish Func类型的variable之后, 访问/debug/vars将会调用该variable的String方法, 而String方法中会调用Func函数, 因此可动态的获取值.  
cmdline和memstats是对Func的典型应用:
```go
func cmdline() interface{} {
	return os.Args
}

func memstats() interface{} {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	return *stats
}
```
在初始化时, 将cmdline和memstats函数强制转换为Func, 之后再Publish: `Publish("cmdline", Func(cmdline))`


### 示例
统计访问/debug/vars的次数.
```go
func expvarMain() {
	expvar.Publish("visitCnt", expvar.Func(visitCnt))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(1)
	}
}

var times = 0

func visitCnt() interface{} {
	times++
	return times
}
```
启动程序后, 可在浏览器中访问`http://localhost:8080/debug/vars`, 刷新可看到visitCnt值的变化.


links
-----
+ [目录](../golang)
+ 上一节: [Pkg sort](Pkg-sort.md)
+ 下一节: [Pkg bytes](Pkg-bytes.md)
