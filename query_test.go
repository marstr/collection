package collection

import "fmt"
import "sync"
import "testing"
import "runtime"
import "time"

func ExampleEnumerator_Any() {
	empty := AsEnumerator()
	if empty.Any() {
		fmt.Println("Empty had some")
	} else {
		fmt.Println("Empty had none")
	}

	populated := AsEnumerator("str1")
	if populated.Any() {
		fmt.Println("Populated had some")
	} else {
		fmt.Println("Populated had none")
	}
	// Output:
	// Empty had none
	// Populated had some
}

func ExampleEnumerator_Count() {
	subject := AsEnumerator("str1", "str1", "str2")
	count1 := subject.Count(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerator_CountAll() {
	subject := AsEnumerator('a', 'b', 'c', 'd', 'e')
	fmt.Println(subject.CountAll())
	// Ouput: 5
}

func ExampleEnumerator_ElementAt() {
	subject := AsEnumerator(1, 2, 3, 5, 8)
	fmt.Print(subject.ElementAt(2))
	// Output: 3
}

func ExampleEnumerator_Last() {
	subject := AsEnumerator(1, 2, 3)
	fmt.Print(subject.Last())
	//Output: 3
}

func ExampleEnumerator_Merge() {
	a := AsEnumerator(1, 2, 4)
	b := AsEnumerator(8, 16, 32)
	c := a.Merge(b)
	sum := 0
	for x := range c {
		sum += x.(int)
	}
	fmt.Println(sum)
	// Output: 63
}

func ExampleMerge() {
	a := AsEnumerator(1, 2, 4)
	b := AsEnumerator(8, 16, 32)
	c := Merge(a, b)
	sum := 0
	for x := range c {
		sum += x.(int)
	}
	fmt.Println(sum)
	// Output: 63
}

func ExampleEnumerator_ParallelSelect() {
	a := AsEnumerator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)

	multiplied := a.ParallelSelect(func(num interface{}) interface{} {
		time.Sleep(50 * time.Millisecond)
		return num.(int) * 5
	})

	sum := 0
	for entry := range multiplied {
		sum = sum + entry.(int)
	}
	fmt.Print(sum)
	// Output: 525
}

func ExampleEnumerator_Reverse() {
	a := AsEnumerator(1, 2, 3)
	fmt.Println(a.Reverse().ToSlice())
	// Output: [3 2 1]
}

func ExampleEnumerator_Select() {
	subject := AsEnumerator('a', 'b', 'c')
	const offset = 'a' - 1
	results := subject.Select(func(a interface{}) interface{} {
		return a.(rune) - offset
	})

	fmt.Println(results.ToSlice())
	// Output: [1 2 3]
}

func ExampleEnumerator_SelectMany() {

	type BrewHouse struct {
		Name  string
		Beers []interface{}
	}

	breweries := AsEnumerator(
		BrewHouse{
			"Mac & Jacks",
			[]interface{}{
				"African Amber",
				"Ibis IPA",
			},
		},
		BrewHouse{
			"Post Doc",
			[]interface{}{
				"Prereq Pale",
			},
		},
		BrewHouse{
			"Resonate",
			[]interface{}{
				"Comfortably Numb IPA",
				"Lithium Altbier",
			},
		},
		BrewHouse{
			"Triplehorn",
			[]interface{}{
				"Samson",
				"Pepper Belly",
			},
		},
	)

	beers := breweries.SelectMany(func(brewer interface{}) Enumerator {
		return AsEnumerator(brewer.(BrewHouse).Beers...)
	})

	for beer := range beers {
		fmt.Println(beer)
	}

	// Output:
	// African Amber
	// Ibis IPA
	// Prereq Pale
	// Comfortably Numb IPA
	// Lithium Altbier
	// Samson
	// Pepper Belly
}

func ExampleEnumerator_Single() {
	a := AsEnumerator(1, 2, 3)
	b := AsEnumerator(4)
	if val, err := a.Single(); err == nil {
		fmt.Println(val)
	}

	if val, err := b.Single(); err == nil {
		fmt.Println(val)
	}
	// Output: 4
}

func ExampleEnumerator_Skip() {
	subject := AsEnumerator(1, 2, 3, 4, 5, 6, 7)
	subject = subject.Skip(5)
	for entry := range subject {
		fmt.Println(entry)
	}
	// Output:
	// 6
	// 7
}

func ExampleEnumerator_Split() {
	a := AsEnumerator(1, 2, 4, 8, 16)
	left, right := a.Split(Identity)

	var wg sync.WaitGroup
	wg.Add(2)

	leftSum := 0
	go func() {
		for x := range left {
			leftSum += x.(int)
		}
		wg.Done()
	}()

	rightSum := 0
	go func() {
		for y := range right {
			rightSum += y.(int)
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Print(leftSum + rightSum)
	// Output: 31
}

func ExampleEnumerator_SplitN() {
	rawNums := make([]interface{}, 1000, 1000)
	for i := 0; i < len(rawNums); i++ {
		rawNums[i] = i + 1
	}
	nums := AsEnumerator(rawNums...)

	workers := nums.SplitN(func(num interface{}) interface{} {
		return num.(int) * 3
	}, 8)

	sum := 0
	results := Merge(workers...)
	for num := range results {
		sum += num.(int)
	}
	fmt.Println(sum)
	// Output: 1501500
}

func ExampleEnumerator_Take() {
	subject := AsEnumerator(1, 2, 3, 4, 5, 6)
	subject = subject.Skip(2).Take(3)
	for entry := range subject {
		fmt.Println(entry)
	}
	// Output:
	// 3
	// 4
	// 5
}

func ExampleEnumerator_TakeWhile() {
	subject := AsEnumerator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	subject = subject.TakeWhile(func(x interface{}, n uint) bool {
		return x.(int) < 6
	})
	for entry := range subject {
		fmt.Println(entry)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleEnumerator_Tee() {
	base := AsEnumerator(1, 2, 4)
	left, right := base.Tee()
	var wg sync.WaitGroup
	wg.Add(2)

	product := 1
	go func() {
		for x := range left {
			product *= x.(int)
		}
		wg.Done()
	}()

	sum := 0
	go func() {
		for x := range right {
			sum += x.(int)
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("Sum: %d\n", sum)
	fmt.Printf("Product: %d\n", product)
	// Output:
	// Sum: 7
	// Product: 8
}

func ExampleEnumerator_UCount() {
	subject := AsEnumerator("str1", "str1", "str2")
	count1 := subject.UCount(func(a interface{}) bool {
		return a == "str1"
	})
	fmt.Println(count1)
	// Output: 2
}

func ExampleEnumerator_UCountAll() {
	subject := AsEnumerator('a', 2, "str1")
	fmt.Println(subject.UCountAll())
	// Output: 3
}

func ExampleEnumerator_Where() {
	subject := AsEnumerator(1, 2, 3, 5, 8, 13, 21, 34)
	results := subject.Where(func(a interface{}) bool { return a.(int) > 8 }).ToSlice()
	fmt.Println(results)
	// Output: [13 21 34]
}

func BenchmarkEnumerator_SplitN(b *testing.B) {
	nums := AsEnumerable(getInitializedSequentialArray()...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := Merge(nums.Enumerate().SplitN(sleepIdentity, uint(runtime.NumCPU()))...)
		for range results {
			// Intentionally Left Blank
		}
	}
}

func BenchmarkEnumerator_Sum(b *testing.B) {
	nums := AsEnumerable(getInitializedSequentialArray()...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range nums.Enumerate().Select(sleepIdentity) {
			// Intentionally Left Blank
		}
	}
}

func sleepIdentity(num interface{}) interface{} {
	time.Sleep(2 * time.Millisecond)
	return Identity(num)
}

func getInitializedSequentialArray() []interface{} {

	rawNums := make([]interface{}, 1000, 1000)
	for i := range rawNums {
		rawNums[i] = i + 1
	}
	return rawNums
}
