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
在单个goroutine中, happens before顺序由程序语句的书写顺序决定.  
*允许*对变量v的read操作r观察到对v的write结果w, 需要满足:
- r不happens before w.
- 在w之后, r之前, 没有额外的write操作.  
请注意, *允许*不代表*一定*, 如果要*确保*对变量v的read操作r观察到对v的write结果w, 需要满足:
- w happens before r.
- 任何对v的write操作必须happens before w, 或者happens after r.  
这对条件比上对条件要求更严格一些: 












links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
