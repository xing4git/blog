Pkg bytes
----

Package bytes实现了一系列操作byte slice方法.

### Compare
比较byte slice a和b, 如果`a == b`返回0, `a < b`返回-1, `a > b`返回1.
```go
func Compare(a, b []byte) int {
	// m取len(a)和len(b)的较小者
	m := len(a)
	if m > len(b) {
		m = len(b)
	}

	// 比较[0, m]之间的byte值
	for i, ac := range a[0:m] {
		bc := b[i]
		switch {
		case ac > bc:
			return 1
		case ac < bc:
			return -1
		}
	}

	// 如果[0, m]之间的byte值都相等, 则判断长度
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}

	// [0,m]之间的byte值都相等, 且a和b的长度也相等, 因此a == b
	return 0
}
```

### Equal
判断2个byte slice是否相等. 此方法不是由Go语言实现的.

### Count
计算s中sep出现的次数.
```go
func Count(s, sep []byte) int {
	n := len(sep)

	// n == 0 和 n > len(s) 这2种情况单独考虑, 直接返回结果
	if n == 0 {
		return utf8.RuneCount(s) + 1
	}
	if n > len(s) {
		return 0
	}

	count := 0
	c := sep[0]
	i := 0
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
			o := IndexByte(t[i:], c)
			// o < 0 时, 说明t[i:]不包含c, 即t[i:]中不可能包含sep了, 所以跳出循环
			if o < 0 {
				break
			}
			// i += o后, t[i] == c
			i += o
		}

		// 此时t[i] == c, 所以只要n == 1 或 Equal(s[i:i+n], sep) == true 都说明有新的sep出现
		if n == 1 || Equal(s[i:i+n], sep) {
			count++
			// 这n个字节相等, 下次从i+n处开始判断
			i += n
			continue
		}

		// 从下个字节开始判断
		i++
	}

	return count
}
```

### Index
计算s中首个sep的索引, -1表示s中不包含sep.
```go
func Index(s, sep []byte) int {
	n := len(sep)
	if n == 0 {
		return 0
	}
	if n > len(s) {
		return -1
	}

	// sep只包含一个byte时, 直接调用IndexByte即可
	c := sep[0]
	if n == 1 {
		return IndexByte(s, c)
	}

	// 这部分代码与Count中类似
	i := 0
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
			o := IndexByte(t[i:], c)
			if o < 0 {
				break
			}
			i += o
		}
		if Equal(s[i:i+n], sep) {
			return i
		}
		i++
	}
	return -1
}
```

### 

links
-----
+ [目录](../golang)
+ 上一节: [Pkg expvar](Pkg-expvar.md)
