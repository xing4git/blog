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
	// 头指针, 尾指针
	front, back *Element
	len         int
}
```

### Next和Prev
```go
// 获取e的后一个Element. nil表示e处于链表的末尾
func (e *Element) Next() *Element { return e.next }
// 获取e的前一个Element. nil表示e处于链表的开头
func (e *Element) Prev() *Element { return e.prev }
```

### Init和New
Init重置list为初始状态
```go
func (l *List) Init() *List {
	// 将list的各个字段设置为对应的零值
	// new(List)返回的*List已经是可以使用的List的, 并不需要调用一次Init后再使用.
	l.front = nil
	l.back = nil
	l.len = 0
	return l
}
// 创建List
func New() *List { return new(List) }
```

### Front和Back
```go
// 返回list中的第一个Element
func (l *List) Front() *Element { return l.front }
// 返回list中的最后一个Element
func (l *List) Back() *Element { return l.back }
```

### Remove
将e从list中删除, 并返回e包含的Value
```go
func (l *List) Remove(e *Element) interface{} {
	l.remove(e)
	e.list = nil // do what remove does not
	return e.Value
}
// 执行实际的删除操作
func (l *List) remove(e *Element) {
	// 判断e是否属于list
	if e.list != l {
		return
	}
	
	if e.prev == nil {
		// 如果e处于链表的开头, 将链表的头指针指向e的下一个Element
		l.front = e.next
	} else {
		// 将e的前一个Element的next指针指向e的后一个Element
		e.prev.next = e.next
	}
	if e.next == nil {
		// 如果e处于链表的末尾, 将链表的尾指针指向e的前一个Element
		l.back = e.prev
	} else {
		// 将e的后一个Element的prev指针指向e的前一个Element
		e.next.prev = e.prev
	}

	e.prev = nil
	e.next = nil
	l.len--
}
```

### insertBefore和insertAfter
```go
// 在mark之前插入e
func (l *List) insertBefore(e *Element, mark *Element) {
	if mark.prev == nil {
		// 头指针指向e
		l.front = e
	} else {
		// 将mark的前一个Element的next指针指向e
		mark.prev.next = e
	}
	// 调整e和mark的前后指针
	e.prev = mark.prev
	mark.prev = e
	e.next = mark
	l.len++
}
// 在mark之后插入e
func (l *List) insertAfter(e *Element, mark *Element) {
	if mark.next == nil {
		// 尾指针指向e
		l.back = e
	} else {
		// 将mark的后一个Element的prev指针指向e
		mark.next.prev = e
	}
	// 调整e和mark的前后指针
	e.next = mark.next
	mark.next = e
	e.prev = mark
	l.len++
}
```

### insertFront和insertBack
```go
// 在链表的开头插入e
func (l *List) insertFront(e *Element) {
	if l.front == nil {
		// 空链表单独处理
		l.front, l.back = e, e
		e.prev, e.next = nil, nil
		l.len = 1
		return
	}
	l.insertBefore(e, l.front)
}
// 在链表的末尾插入e
func (l *List) insertBack(e *Element) {
	if l.back == nil {
		// 空链表单独处理
		l.front, l.back = e, e
		e.prev, e.next = nil, nil
		l.len = 1
		return
	}
	l.insertAfter(e, l.back)
}
```

### PushFront和PushBack
```go
// 将value包装成Element, 并将Element插入链表开头
func (l *List) PushFront(value interface{}) *Element {
	e := &Element{nil, nil, l, value}
	l.insertFront(e)
	return e
}
// 将value包装成Element, 并将Element插入链表末尾
func (l *List) PushBack(value interface{}) *Element {
	e := &Element{nil, nil, l, value}
	l.insertBack(e)
	return e
}
```

### InsertBefore和InsertAfter
```go
// 将value包装成Element, 并将Element插入mark之前
func (l *List) InsertBefore(value interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	e := &Element{nil, nil, l, value}
	l.insertBefore(e, mark)
	return e
}
// 将value包装成Element, 并将Element插入mark之后
func (l *List) InsertAfter(value interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	e := &Element{nil, nil, l, value}
	l.insertAfter(e, mark)
	return e
}
```

### MoveToFront和MoveToBack
```go
// 将e移动到链表的开头
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.front == e {
		return
	}
	// 先将e删除, 然后在链表开头插入
	l.remove(e)
	l.insertFront(e)
}

// 将e移动到链表的末尾
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.back == e {
		return
	}
	// 先将e删除, 然后在链表末尾插入
	l.remove(e)
	l.insertBack(e)
}
```

### PushFrontList和PushBackList
```go
// 将ol整体插入链表的开头
func (l *List) PushFrontList(ol *List) {
	first := ol.Front()
	// 从ol的尾部开始, 依次将ol的每个value插入链表的开头
	for e := ol.Back(); e != nil; e = e.Prev() {
		l.PushFront(e.Value)
		if e == first {
			break
		}
	}
}
// 将ol整体插入链表的末尾
func (l *List) PushBackList(ol *List) {
	last := ol.Back()
	// 从ol的头部开始, 依次将ol的每个value插入链表的末尾
	for e := ol.Front(); e != nil; e = e.Next() {
		l.PushBack(e.Value)
		if e == last {
			break
		}
	}
}
```

### Len
返回链表的长度
```go
func (l *List) Len() int { return l.len }
```



links
-----
+ [目录](../golang)
+ 上一节: [Pkg ioutil](Pkg-ioutil.md)
+ 下一节: [Go Memory Model](Go-Memory-Model.md)
