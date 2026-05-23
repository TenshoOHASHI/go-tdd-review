package interface_test

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {

	rectangle := Rectangle{10.0, 10.10}
	got := Perimeter(rectangle.Width, rectangle.Height)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	// インターフェースに満たす、各面積を呼び出す
	// checkArea := func(t *testing.T, shape Shape, want float64) {
	// 	t.Helper()
	// 	got := shape.Area()

	// 	if got != want {
	// 		// gを使用すると、10進数で表示
	// 		t.Errorf("got %g, want %g", got, want)
	// 	}
	// }

	// 構造体を持つスライス＝スライスに複数の構造体を用意
	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		// 一番外側がスライス全体の中身＝[elem1, elem2]
		//　内側が構造体を初期化する値＝shape: value, want:value
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		// 2つめの初期化
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
			}
		})
	}
}

type Shape interface {
	Area() float64
}

// 矩形のデータフィールド
type Rectangle struct {
	Width  float64
	Height float64
}

// 円のデータフィールド
type Circle struct {
	Radius float64
}

// 　三角形のフォールド
type Triangle struct {
	Base   float64
	Height float64
}

// 周囲の長さ
func Perimeter(width, height float64) float64 {
	return 2 * (width + height)
}

// 面積
func Area(width float64, height float64) float64 {
	return width * height
}

// レシーバーメソッド
// インスタンスに対してメソッドを呼び出す
// データフィールドの参照を受け取る
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// 矩形の面積
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 三角形の面積
func (t Triangle) Area() float64 {
	return (t.Base * t.Height) / 2
}
