package main

import (
	"fmt"
	"regexp"

	"github.com/x0xO/hg"
)

func main() {
	// strings
	str := hg.NewHString() // declaration and assignation
	fmt.Println(str.Random(9))

	str = "test"

	fmt.Println(hg.NewHString("12").HInt())
	fmt.Println(str.Hash().MD5())
	fmt.Println(str.Encode().GzFlate().Decode().GzFlate())

	var str2 hg.HString = "rest" // declaration and assignation

	fmt.Println(str2)

	a := hg.NewHString("abc")
	b := hg.NewHString("bbb")
	c := hg.NewHString("ccc")
	d := hg.NewHString("ddd")
	e := hg.NewHString("eee")

	str3 := a.ReplaceAll("a", "zzz").ToUpper().Fields().Random().Split("")[0].ToLower().String()

	fmt.Println(str3)

	// ints
	n := hg.NewHInt(52452356235) // declaration and assignation

	fmt.Printf("%v\n", n.Bytes())

	fmt.Println(n.Hash().MD5())
	fmt.Println(n.Hash().SHA1())
	fmt.Println(n.Hash().SHA256())

	fmt.Println(n.ToBinary())
	fmt.Println(n.HString())

	rn := hg.NewHInt(10).Random()
	fmt.Println("random number: ", rn)

	rrn := hg.NewHInt().RandomRange(10, 100)
	fmt.Println("random range number: ", rrn)

	var n2 hg.HInt = 321

	fmt.Println(n2) // declaration and assignation

	fmt.Println(n.Add(n2).Mul(3))

	// floats

	fl := hg.NewHFloat(12.5456)
	fmt.Println(fl.Round()) // 13

	// slices
	sss := hg.HSlice[int]{}
	fmt.Println(sss.Chunks(2))

	sl := hg.NewHSlice[hg.HString]().Append(a, b, c, d, e) // declaration and assignation

	sl.Shuffle()

	fmt.Println(sl.Get(-1) == "eee")
	fmt.Println(sl.Get(1) == "bbb")
	fmt.Println(sl.Get(-2) == "ddd")

	fmt.Println(sl.Map(hg.HString.ToUpper))

	fmt.Println(sl.Permutations())

	counter := sl.Append(sl...).Append("ddd").Counter()
	fmt.Println(counter) // Output: HMapOrd[bbb:2, eee:2, ccc:2, abc:2, ddd:3]

	counter.ForEach(func(k any, v int) { fmt.Println(k.(hg.HString).ToTitle(), ":", v) })

	sl.ForEach(func(v hg.HString) { fmt.Println(v) })

	sl = sl.Unique().Reverse().Filter(func(s hg.HString) bool { return s != "bbb" })

	fmt.Println(sl.Random())

	sl1 := hg.HSliceOf(1, 2, 3, 4, 5) // declaration and assignation

	fmt.Println(sl1.Reduce(func(index, value int) int { return index + value }, 0)) // 15

	sl3 := hg.HSlice[hg.HString]{} // declaration and assignation
	sl3 = sl3.Append("aaaaa", "bbbbb")

	fmt.Println(sl3.ToHMapHashed())

	fmt.Println(sl3.Last().Count("b")) // 5

	sl4 := hg.HSliceOf([]string{"root", "toor"}...).Random()
	fmt.Println(hg.NewHString(sl4).ToUpper())

	fmt.Println(sl3.Map(func(s hg.HString) hg.HString { return s + "MAPMAPMAP" }))
	fmt.Println(sl3.MapParallel(func(s hg.HString) hg.HString { return s + "MAPMAPMAP" }))

	empsl := hg.NewHSlice[hg.HString]()

	fmt.Println(empsl.Empty())

	// maps
	m1 := hg.HMapFromMap(map[int]string{1: "root", 22: "toor"}) // declaration and assignation
	fmt.Println(m1.Values())

	m2 := hg.NewHMap[int, string]() // declaration and assignation

	m2[99] = "AAA"
	m2[88] = "BBB"
	m2.Set(77, "CCC")

	fmt.Println(m2.Delete(99))
	fmt.Println(m2.Keys())

	fmt.Println(m2)
	fmt.Println(m2.ToMap())

	fmt.Println(m2.Invert().Values().Get(0))        // return int type
	fmt.Println(m2.Invert().Keys().Get(0).(string)) // return any type, need assert to type

	m3 := hg.HMap[string, string]{"test": "rest"} // declaration and assignation
	fmt.Println(m3.Contains("test"))

	slp := hg.NewHSlice[int](2049 * 51).Fill(22)
	random := hg.NewHInt(99).Random().Int()

	slp = slp.MapParallel(func(i int) int { return i * random })

	slp = slp.FilterParallel(func(i int) bool { return i%2 == 0 })
	fmt.Println(slp.ReduceParallel(func(index, value int) int { return index + value }, 0))

	ub := hg.NewHBytes([]byte("abcdef\u0301\u031dg"))
	fmt.Println(ub.NormalizeNFC().Reverse())

	us := hg.NewHString("abcde丂g")
	fmt.Println(us.Reverse())

	l := hg.HString("hello")
	fmt.Println(l.Similarity("world"))

	hbs := hg.HBytes([]byte("Hello, 世界!"))
	reversed := hbs.Reverse() // "!界世 ,olleH"

	fmt.Println(reversed)

	hbs = hg.HBytes([]byte("hello, world!"))
	replaced := hbs.Replace([]byte("l"), []byte("L"), 2) // "heLLo, world!"

	fmt.Println(replaced)

	hs1 := hg.HString("kitten")
	hs2 := hg.HString("sitting")
	similarity := hs1.Similarity(hs2) // hg.HFloat(57.14285714285714)

	fmt.Println(similarity)

	fmt.Println(hg.NewHString("&aacute;").Decode().HTML())

	to := hg.HString("Hello, 世界!")

	fmt.Println(to.Encode().Binary().Chunks(8).Join(" "))
	fmt.Println(to.Encode().Binary().Decode().Binary())

	fmt.Println(to.Encode().Hex())
	fmt.Println(to.Encode().Hex().Decode().Hex())

	toi := hg.HInt(1234567890)

	fmt.Println(toi.ToBinary())
	fmt.Println(toi.ToOctal())
	fmt.Println(toi.ToHex())

	ascii := hg.HString("💛💚💙💜")
	fmt.Println(ascii.IsASCII())

	reg := hg.NewHString("some text")
	fmt.Println(reg.ContainsRegexp(regexp.MustCompile(`\w+`)))

	fmt.Println(hg.HString("example.com").EndsWith(".com", ".net"))

	fmt.Println(hg.NewHString("Hello").Format("%s world"))
}
