Pkg bytes
----

Package bytes实现了一系列操作byte slice的方法.

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
s中首次出现sep序列的索引, -1表示s中不包含sep.
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

### LastIndex
s中最后出现sep序列的索引, -1同样表示s中不包含sep.
```go
func LastIndex(s, sep []byte) int {
	n := len(sep)
	if n == 0 {
		return len(s)
	}

	// 这里的代码比Index中的代码更容易理解
	c := sep[0]
	for i := len(s) - n; i >= 0; i-- {
		if s[i] == c && (n == 1 || Equal(s[i:i+n], sep)) {
			return i
		}
	}

	return -1
}
```


### Contains
判断b中是否包含指定的subslice.
```go
func Contains(b, subslice []byte) bool {
	return Index(b, subslice) != -1
}
```

### IndexRune
将s当做utf8编码的unicode字符序列, 返回r在s中的索引.
```go
func IndexRune(s []byte, r rune) int {
	for i := 0; i < len(s); {
		// 返回s[i:]中的第一个rune, size为该rune所占的字节数
		r1, size := utf8.DecodeRune(s[i:])
		if r == r1 {
			return i
		}
		i += size
	}
	return -1
}
```

### IndexAny
chars中任意一个字符在s中首次出现的索引.
```go
func IndexAny(s []byte, chars string) int {
	if len(chars) > 0 {
		var r rune
		var width int

		for i := 0; i < len(s); i += width {
			// r为s[i:]中的第一个rune, width为该rune所占的字节数
			r = rune(s[i])
			if r < utf8.RuneSelf {
				width = 1
			} else {
				r, width = utf8.DecodeRune(s[i:])
			}

			// 如果r是chars其中一员, 直接返回
			for _, ch := range chars {
				if r == ch {
					return i
				}
			}
		}
	}
	return -1
}
```

### LastIndexAny
chars中任意一个字符在s中最后出现的索引.
```go
func LastIndexAny(s []byte, chars string) int {
	if len(chars) > 0 {
		for i := len(s); i > 0; {
			r, size := utf8.DecodeLastRune(s[0:i])
			i -= size
			for _, ch := range chars {
				if r == ch {
					return i
				}
			}
		}
	}
	return -1
}
```

### Split
bytes包中包含多个Split函数: SplitN, SplitAfterN, Split, SplitAfter等. 这些函数都是对genSplit的包装.
```go
// 将s以sep分割. sepSave表示子byte数组中包含多少个sep中的字节数.
// n表示最多包含多少个子byte数组:
// n > 0 取n个子byte数组, n == 0 返回nil, n < 0 时返回所有子byte数组
func genSplit(s, sep []byte, sepSave, n int) [][]byte {
	// 根据n的值确定要返回多少个byte数组
	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		return explode(s, n)
	}
	if n < 0 {
		n = Count(s, sep) + 1
	}

	c := sep[0]
	start := 0
	a := make([][]byte, n)
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || Equal(s[i:i+len(sep)], sep)) {
			// 找到其中一个, 从中可以看出sepSave的作用
			a[na] = s[start : i+sepSave]
			na++
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	// 剩余的数据都归类到最后一个子byte数组
	a[na] = s[start:]

	return a[0 : na+1]
}

// 将s以sep作为分界点分割为最多n个子byte数组, 每个子byte不数组包含sep
func SplitN(s, sep []byte, n int) [][]byte { return genSplit(s, sep, 0, n) }

// 将s以sep作为分界点分割为最多n个子byte数组, 每个子byte数组包含sep
func SplitAfterN(s, sep []byte, n int) [][]byte {
	return genSplit(s, sep, len(sep), n)
}

// 将s以sep作为分界点分割为若干个子byte数组, 不限制子byte数组的个数. 每个子byte不数组包含sep
func Split(s, sep []byte) [][]byte { return genSplit(s, sep, 0, -1) }

// 将s以sep作为分界点分割为若干个子byte数组, 不限制子byte数组的个数. 每个子byte数组包含sep
func SplitAfter(s, sep []byte) [][]byte {
	return genSplit(s, sep, len(sep), -1)
}
```

### Fields
将s以空白字符作为分界点分割为若干个子byte数组.
```go
func Fields(s []byte) [][]byte {
	return FieldsFunc(s, unicode.IsSpace)
}

// 对于s中的每个rune, 都调用f函数. 如果f函数返回true, 则将其作为分界点.
func FieldsFunc(s []byte, f func(rune) bool) [][]byte {
	// 计算得到子byte数组的长度
	n := 0
	inField := false
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		wasInField := inField
		inField = !f(r)
		if inField && !wasInField {
			n++
		}
		i += size
	}

	a := make([][]byte, n)
	na := 0
	fieldStart := -1
	for i := 0; i <= len(s) && na < n; {
		r, size := utf8.DecodeRune(s[i:])
		if fieldStart < 0 && size > 0 && !f(r) {
			fieldStart = i
			i += size
			continue
		}
		if fieldStart >= 0 && (size == 0 || f(r)) {
			a[na] = s[fieldStart:i]
			na++
			fieldStart = -1
		}
		if size == 0 {
			break
		}
		i += size
	}
	return a[0:na]
}
```

### Join
将a的每一个item通过sep连接在一起, 形成一个大的byte数组.
```go
func Join(a [][]byte, sep []byte) []byte {
	// 处理特殊情况
	if len(a) == 0 {
		return []byte{}
	}
	if len(a) == 1 {
		return a[0]
	}

	// 计算返回结果的len
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	// copy
	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}

	return b
}
```

### HasPrefix, HasSuffix
代码很简单.
```go
// 是否以prefix开端
func HasPrefix(s, prefix []byte) bool {
	return len(s) >= len(prefix) && Equal(s[0:len(prefix)], prefix)
}
// 是否以suffix结尾
func HasSuffix(s, suffix []byte) bool {
	return len(s) >= len(suffix) && Equal(s[len(s)-len(suffix):], suffix)
}
```



links
-----
+ [目录](../golang)
+ 上一节: [Pkg expvar](Pkg-expvar.md)
