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
在单个goroutine中, happens before顺序由程序语句的书写顺序指定. 对变量v的read能够获知对v的write结果, 需要满足:
- read不happens before write.
- 在write之后, read之前, 没有额外的write操作.



links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
