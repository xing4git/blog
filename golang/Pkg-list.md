Pkg list
----

list包提供双向链表的实现

### List和Element
```go
// List中的单个元素
type Element struct {
	// 前后指针
	next, prev *Element
	// The list to which this element belongs.
	list *List
	// 包含的实际值
	Value interface{}
}

type List struct {
	// 首尾指针
	front, back *Element
	len         int
}
```


links
-----
+ [目录](../golang)
+ 上一节: [Pkg ioutil](Pkg-ioutil.md)
