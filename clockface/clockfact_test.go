package clockface_test

import (
	"math"
	"test/clockface"
	"testing"
	"time"
)

func TestSecondHandAtMidnight(t *testing.T) {

	// 24時、秒針が一番上
	// Y軸は上方向でマイナス
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 - 90}

	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

// func TestSecondHandAt30Seconds(t *testing.T) {
// 	// 秒針が一番下に傾く、ためY軸方向にプラス９０の長さを追加
// 	// 360/60 * 30 = 180度
// 	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

// 	want := clockface.Point{X: 150, Y: 150 + 90}
// 	got := clockface.SecondHand(tm)

// 	if got != want {
// 		t.Errorf("Got %v, wanted %v", got, want)
// 	}
// }

func TestSecondsInRadians(t *testing.T) {
	cases := []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 0, 30), math.Pi},           // 180度
		{simpleTime(0, 0, 0), 0},                  // 0度
		{simpleTime(0, 0, 45), (math.Pi / 2) * 3}, // 270度、
		{simpleTime(0, 0, 7), (math.Pi / 30) * 7},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondsInRadians(c.time)
			if got != c.angle {
				t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
			}
		})
	}
}

type Point struct {
	X float64
	Y float64
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

// math.Piを因数分解すると、１０進数の値が変わってしまう。
func secondsInRadians(t time.Time) float64 {
	// return (2 * (math.Pi) / float64(60) * float64(t.Second()))
	return (math.Pi) / (30 / float64(t.Second()))
}

func secondHandPoint(t time.Time) Point {

	// 角度を取得
	angle := secondsInRadians(t)
	// 角度を使って座標を計算
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Point{X: x, Y: y}
}

func TestSecondHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	// a-bはその差の絶対値、必ずプラスにし、小数点以下７桁ならそれを返す
	return math.Abs(a-b) < equalityThreshold

}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}

const secondHandLength = 90
const clockCentreX = 150
const clockCentreY = 150

func SecondHand(t time.Time) Point {
	p := secondHandPoint(t)
	p = Point{p.X * secondHandLength, p.Y * secondHandLength} // scale
	p = Point{p.X, -p.Y}                                      // flip
	p = Point{p.X + clockCentreX, p.Y + clockCentreY}         //translate
	return p
}

func TestSecondHandAt30Seconds(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

	// y軸をスケール
	want := Point{X: 150, Y: 150 + 90}
	got := SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
