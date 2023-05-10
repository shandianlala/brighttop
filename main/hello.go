package main

import "fmt"

// 运行方式
// go run hello.go
// ./hello

type List []int

func (l List) Len() int {
	return len(l)
}

func (l *List) Append(val int) {
	*l = append(*l, val)
}

type Appender interface {
	Append(int)
}

func CountInto(a Appender, start, end int) {
	for i := start; i <= end; i++ {
		a.Append(i)
	}
}

type Lener interface {
	Len() int
}

func LongEnough(l Lener) bool {
	return l.Len()*10 > 42
}

func main() {
	// A bare value
	var lst List
	// compiler error:
	// cannot use lst (type List) as type Appender in argument to CountInto:
	//       List does not implement Appender (Append method has pointer receiver)
	//CountInto(lst, 1, 10)
	if LongEnough(lst) { // VALID: Identical receiver type
		fmt.Printf("- lst is long enough\n")
	}

	// A pointer value
	plst := new(List)
	CountInto(plst, 1, 10) // VALID: Identical receiver type
	if LongEnough(plst) {
		// VALID: a *List can be dereferenced for the receiver
		fmt.Printf("- plst is long enough\n")
	}

	var p *int                // 定义了指针变量p，指针变量p还未初始化
	fmt.Println("1. ", p, &p) //  1.  <nil> 0xc000124088
	p = new(int)
	fmt.Println("2. ", p, &p, *p) //  2.  0xc000128c00 0xc000124088 0
	*p = 1
	fmt.Println("3. ", p, &p, *p) //  3.  0xc000128c00 0xc000124088 1

	// first
	/*second*/
	fmt.Println("Hello, World!")
	// %d 表示整型数字，%s 表示字符串
	var stockcode = 123
	var enddate = "2020-12-31"
	var url = "Code=%d&endDate=%s"
	var target_url = fmt.Sprintf(url, stockcode, enddate)
	fmt.Println(target_url)

	var a string = "Runoob"
	fmt.Println(a)
	var b, c int = 1, 2
	fmt.Println(b, c)
	notInitVarValue()

	f := "Runoob = var f string = \"Runoob\"" // var f string = "Runoob"
	fmt.Println(f)

	//这种不带声明格式的只能在函数体中出现
	//g, h := 123, "hello"

	// 一个局部变量却没有在相同的代码块中使用它，同样会得到编译错误
	// 但是全局变量是允许声明但不使用的

	// 空白标识符 _ 也被用于抛弃值
	_, numb, strs := numbers() //只获取函数返回值的后两个
	fmt.Println(numb, strs)

	var a1 int = 4
	var ptr *int
	/* 运算符实例 */
	fmt.Printf("第 1 行 - a 变量类型为 = %T\n", a)
	/*  & 和 * 运算符实例 */
	ptr = &a1 /* 'ptr' 包含了 'a' 变量的地址 */
	fmt.Printf("a 的值为  %d\n", a1)
	fmt.Println(ptr)
	fmt.Printf("ptr 变量的地址是: %x\n", ptr)
	fmt.Printf("*ptr 为 %d\n", *ptr)

	numbers := [6]int{1, 2, 3, 5}
	for i, x := range numbers {
		fmt.Printf("第 %d 位 x 的值 = %d\n", i, x)
	}

	v1, v2 := swap("p1", "p2", 66)
	fmt.Println("v1 = " + v1 + ", v2 = " + v2)
	testmap()

}

func testmap() {
	var countryCapitalMap map[string]string /*创建集合 */
	countryCapitalMap = make(map[string]string)

	/* map插入key - value对,各个国家对应的首都 */
	countryCapitalMap["France"] = "巴黎"
	countryCapitalMap["Italy"] = "罗马"
	countryCapitalMap["Japan"] = "东京"
	countryCapitalMap["India "] = "新德里"

	/*使用键输出地图值 */
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap[country])
	}
	for key, value := range countryCapitalMap {
		fmt.Printf("key -> %s, value -> %s\n", key, value)
	}

	/*查看元素在集合中是否存在 */
	capital, ok := countryCapitalMap["American"] /*如果确定是真实的,则存在,否则不存在 */
	/*fmt.Println(capital) */
	/*fmt.Println(ok) */
	if ok {
		fmt.Println("American 的首都是", capital)
	} else {
		fmt.Println("American 的首都不存在")
	}
}

func swap(x, y string, index int) (string, string) {
	return y + "_" + string(rune(index)), x + "_" + string(rune(index))
}

//一个可以返回多个值的函数
func numbers() (int, int, string) {
	a, b, c := 1, 2, "str"
	return a, b, c
}

func notInitVarValue() {
	// 声明一个变量并初始化
	var a = "RUNOOB"
	fmt.Println(a)

	// 没有初始化就为零值
	var b int
	fmt.Println(b)

	// bool 零值为 false
	var c bool
	fmt.Println(c)
}
