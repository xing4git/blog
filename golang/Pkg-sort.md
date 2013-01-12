Pkg sort
----

Package sort为slice或者自定义的集合提供排序功能.

### Interface
Interface定义了排序所需实现的接口.
```go
type Interface interface {
  // 集合中元素的个数
	Len() int
	// 第i个元素和第j个元素的比较
	Less(i, j int) bool
	// 交换第i个元素和第j个元素的位置
	Swap(i, j int)
}
```
如果想要使用sort包提供的排序功能, 相应的类型就必须实现Interface接口.

### Sort
Sort方法实现排序.
```go
func Sort(data Interface) {
	n := data.Len()
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(data, 0, n, maxDepth)
}
```
内部使用堆排序和快速排序结合的算法:
```go
func quickSort(data Interface, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(data, a, b)
	}
}
```

### IsSorted
判断数据是否是升序排列的.
```go
func IsSorted(data Interface) bool {
	n := data.Len()
	for i := n - 1; i > 0; i-- {
	  // 如果出现了降序, 则直接返回false
		if data.Less(i, i-1) {
			return false
		}
	}
	return true
}
```

### IntSlice
为[]int类型实现Interface接口.
```go
type IntSlice []int
func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IntSlice) Sort() { Sort(p) }
```
对于使用者来说, 之需要这样:
```go
ints := []int{1, 9, 3, -1, 18}
// 强制转换
intSlice := sort.IntSlice(ints)
intSlice.Sort()
fmt.Println(ints) // 输出: [-1 1 3 9 18]
```

### Float64Slice, StringSlice
类似的, sort包还为[]float64和[]string实现了Interface接口, 以方便我们的使用.
```go
type Float64Slice []float64
func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return p[i] < p[j] || math.IsNaN(p[i]) && !math.IsNaN(p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float64Slice) Sort() { Sort(p) }

type StringSlice []string
func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Sort() { Sort(p) }
```

### Ints, Float64s, Strings
更为贴心的是, sort包还为我们可以更为简便的排序方法.
```go
func Ints(a []int) { Sort(IntSlice(a)) }
func Float64s(a []float64) { Sort(Float64Slice(a)) }
func Strings(a []string) { Sort(StringSlice(a)) }
```
这样一来, 想要对[]int排序就更简单了:
```go
ints := []int{1, 9, 3, -1, 18}
sort.Ints(ints)
fmt.Println(ints)
```

### IntsAreSorted, Float64sAreSorted, StringsAreSorted
这3个方法也是为了方便调用者使用而存在的.
```go
// 判断[]int是否按升序排列
func IntsAreSorted(a []int) bool { return IsSorted(IntSlice(a)) }
func Float64sAreSorted(a []float64) bool { return IsSorted(Float64Slice(a)) }
func StringsAreSorted(a []string) bool { return IsSorted(StringSlice(a)) }
```

### 为自定义类型实现Interface接口
如果我们自定义了Person类型, 那么如何为[]Person排序呢? 最简单的方法就是让[]Person实现Interface接口.
```go
type Person struct {
	Name string
	Age  int
}
type PersonSlice []Person

func (ps PersonSlice) Len() int {
	return len(ps)
}
func (ps PersonSlice) Less(i, j int) bool {
	return ps[i].Name < ps[j].Name
}
func (ps PersonSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func personSort() {
	ps := []Person{{"xing", 25}, {"min", 30}, {"yong", 29}}
	sort.Sort(PersonSlice(ps))
	fmt.Println(ps) // [{min 30} {xing 25} {yong 29}]
}
```
如果不想按升序排序, 想要降序排列, 之需要稍微更改一下Less方法即可.
```go
func (ps PersonSlice) Less(i, j int) bool {
	return ps[i].Name > ps[j].Name
}
```
如果想要按Age排列也很简单.
```go
func (ps PersonSlice) Less(i, j int) bool {
	return ps[i].Age < ps[j].Age
}
```

### Search
排序总是和查找不分家, 为此, sort包也实现了二分查找法.
```go
func Search(n int, f func(int) bool) int {
	i, j := 0, n
	for i < j {
	  // 取得中间索引
		h := i + (j-i)/2
		if !f(h) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
```
Search方法的参数n和f需要满足条件:
+ 存在索引x, 使得当i属于[0,x]时, f(i) == false. 且
+ 当i属于(x,n)时, f(i) == true.
+ Search方法返回x, the first true index.

### SearchInts, SearchFloat64s, SearchStrings
如果想要在已排好序的[]int, []float64, 或者[]string中查找某个值, sort包已经实现了相应的方法.
```go
// 在a中查找x, 返回x所在的索引. a必须是已经按升序排序的slice.
func SearchInts(a []int, x int) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}
func SearchFloat64s(a []float64, x float64) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}
func SearchStrings(a []string, x string) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}
// 对SearchInts的包装
func (p IntSlice) Search(x int) int { return SearchInts(p, x) }
func (p Float64Slice) Search(x float64) int { return SearchFloat64s(p, x) }
func (p StringSlice) Search(x string) int { return SearchStrings(p, x) }
```
以[]int为例:
```go
ints := []int{1, 9, 3, -1, 18}
// sort first
sort.IntSlice(ints).Sort()
fmt.Println(sort.SearchInts(ints, 9)) // 3
```
当然, 也可以这样:
```go
ints := []int{1, 9, 3, -1, 18}
is := sort.IntSlice(ints)
is.Sort()
fmt.Println(is.Search(9))
```

### 根据Name查找Person
如果想要在[]Person中查找出Name为指定值的Person, 可以这样:
```go
func (ps PersonSlice) Search(name string) Person {
	searchFunc := func(i int) bool {
		return ps[i].Name >= name
	}
	// 查找到the fisrt true index, 即第一个使得ps[i].Name >= name成立的i
	index := sort.Search(len(ps), searchFunc)
	return ps[index]
}
// 使用示例
func personSearch() {
	persons := []Person{{"xing", 25}, {"min", 30}, {"yong", 29}}
	ps := PersonSlice(persons)
	sort.Sort(ps)
	fmt.Println(ps.Search("xing")) // {xing 25}
}
```




links
-----
+ [目录](../golang)
+ 上一节: [Pkg log](Pkg-log.md)
