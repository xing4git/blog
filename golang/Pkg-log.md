Pkg log
----

Package log实现了简单的日志功能. 定义了类型Logger, Logger用于格式化输出日志.
同时, 预定义了名为std的Logger对象, 用于输出日志到stderr.
包中的Fatal方法将在输出日志消息后调用os.Exit(1), 而Panic方法将在输出日志消息后调用panic.



links
-----
+ [目录](../golang)
+ 上一节: [Pkg list](Pkg-list.md)
