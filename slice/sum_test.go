package slice_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		// 固定長の配列
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15
		if want != got {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		got := SumSlice(slice)
		want := 15

		if got != want {
			t.Errorf("get %d want %d given %v", got, want, slice)

		}
	})
}

// 引数は固定長なので、指定する必要がある
func Sum(numbers [5]int) int {
	var total int
	for _, num := range numbers {
		// 加算代入演算子
		total += num
	}

	return total
}

// スライスで値を受け取る
func SumSlice(numbers []int) int {
	sum := 0

	for _, num := range numbers {
		sum += num
	}
	return sum
}

// 複数のスライスを受け取り、新しいスライスを返す
func TestSumAll(t *testing.T) {

	// ここで受け取る引数を固定することで、違った引数が渡されないようにしている
	checkSums := func(t *testing.T, got, want []int) {
		t.Helper()
		// スライスはnilとの比較以外比較ができない。
		// DeepEqualでスライスの中身まで見てくれる
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
	got := SumAllAppend([]int{1, 2}, []int{0, 9})

	want := []int{3, 9}
	fmt.Print("got", got)
	// 同じ型か確認する
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}

	t.Run("空のスライスを渡す", func(t *testing.T) {
		got := SumAll([]int{}, []int{0, 9})
		want := []int{0, 9}
		checkSums(t, got, want)
	})
}

// 1. ... の右側にある []int が、受け取る各要素の型
// 2. 呼び出し側から []int の値を複数受け取れる
// 3. Go がそれらをまとめるために、外側のスライスを作る
// 4. その結果、関数内の numbersToSum は [][]int になる
// ...[]int は「[]int 型の値を複数受け取り、それらを外側のスライスにまとめる」
func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)

	// makeを使うと、データの中身をゼロ値で初期化することができる、そのためインデックス化して値を取得したり、新しい値を代入することができる
	sums := make([]int, lengthOfNumbers)
	// 存在しないインデックスにアクセスすると、エラーになります。
	for i, numbers := range numbersToSum {
		sums[i] = SumSlice(numbers)
	}
	return sums
}

func SumAllAppend(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)

	sums := make([]int, lengthOfNumbers)
	for _, num := range numbersToSum {
		// append は容量が足りない場合、自動でより大きい配列を確保して要素をコピーする。
		// 事前に十分な容量を確保しておくと、その再確保を減らせる。
		sums = append(sums, SumSlice(num))
	}
	return sums
}

func SumAllEmpty(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		}
		sums = append(sums, SumSlice(numbers))
	}
	return sums

}
