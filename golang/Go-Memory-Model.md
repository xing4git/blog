Go Memory Model
----

本文是对[The Go Memory Model](http://golang.org/ref/mem)的翻译和理解.

Introduction
----
Go Memory Model指定在什么条件下, 一个goroutine能够read到另一个goroutine对某个变量的write结果.

Happens Before
----
执行代码时, 编译器和处理器可能优化指令的执行顺序, 这被称为指令重排序. 指令重排序具有这样的原则: 不改变语言规范中定义的行为.  
语言规范中定义了happens before, 用于指定read和write操作的执行顺序. 事件e1 happens before e2, 等价于e2 happens after e1. 如果e1不happens before e2, 也不happens after e2, 那么我们说e1和e2 happen concurrently.  
[允许]()对变量v的read操作r观察到对v的write结果w, 需要满足:
- r不happens before w.
- 不存在额外的write操作w1, happens after w且happens before r.

请注意, [允许]()不代表[一定](), 如果要[确保]()对变量v的read操作r观察到对v的write结果w, 需要满足:
- w happens before r.
- 任何额外的对v的write操作必须happens before w, 或者happens after r.

2对条件看似差不多, 其实第二对条件比第一对条件要求更严格一些: 第二对条件排除了happens concurrently情形, 也就是说r和w不能happens concurrently, 也不能存在额外的write操作happens concurrently with r or w.
在单goroutine环境中, 不可能存在happens concurrently, 所以此时2对条件是等价的. 当多个goroutine可以访问变量v时, 就需要通过同步以确保happens before.


Synchronization
----
+ 初始化  
程序的所有初始化操作都发生在同一个goroutine中. 程序的初始化操作包括: 初始化包级变量和常量, 执行init函数.  
如果包p import包q, 那么q的init函数happens before p的init函数.  
main.main函数happens after所有init函数.  
+ goroutine的创建  
我们知道, goroutine由go语句创建. go语句happens before goroutine的执行. 例如:
```go
var a string
func f() {
	print(a)
}
func hello() {
	a = "hello world"
	go f()
}
```
+ goroutine的结束  
goroutine的结束不happens before于任何事件.
```go
var a string
func hello() {
	go func() { a = "hello" }
	print(a)
}
```
此时不能确保输出的结果一定是"hello".
+ Channel通信  













links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
