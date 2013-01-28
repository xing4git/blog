Index
-----

####Pkg bufio
*2013-01-04*
bufio包提供了带有缓冲功能的Reader和Writer, 应用了`装饰者模式`.  bufio.Reader是对io.Reader的包装, 并且实现了io.Reader接口.  ...[Read More](./Pkg-bufio.md)

####Pkg ioutil
*2013-01-05*
ioutil包提供了一些io操作的便利方法....[Read More](./Pkg-ioutil.md)

####Pkg list
*2013-01-06*
list包提供双向链表的实现...[Read More](./Pkg-list.md)

####Pkg log
*2013-01-11*
Package log实现了简单的日志功能. 定义了类型Logger, 用于输出格式化日志.  同时, 预定义了名为std的Logger对象, 用于输出日志到stderr.  ...[Read More](./Pkg-log.md)

####Pkg sort
*2013-01-12*
Package sort为slice或者自定义的集合提供排序功能....[Read More](./Pkg-sort.md)

####Pkg expvar
*2013-01-13*
Package expvar是用于获取public variables的标准接口. 通过http访问/debug/vars, 即可得到json格式的public variables.  expvar默认publish了2个variable: cmdline和memstats, 分别用于获取程序启动时的命令行参数和当前内存使用情况.  ...[Read More](./Pkg-expvar.md)

####Pkg bytes
*2013-01-26*
Package bytes实现了一系列操作byte slice的方法....[Read More](./Pkg-bytes.md)

