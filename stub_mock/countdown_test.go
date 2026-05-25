package main_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

// 自作の型
type MyWriter struct{}

// からのインスタンスを作成すれば、Sleepメソッドを享受することができる
type DefaultSleeper struct{}

// 複数のテスト結果を保存する
type CountdownOperationsSpy struct {
	Calls []string
}

// 関数インジェクション
type ConfigurableSleeper struct {
	// タイマーを設定
	duration time.Duration
	// 実際に、関数を用意し、タイマーを変数に格納する
	// 関数を値として持つ（メソッドが２つ以上ならインターフェースを使用、１つの場合は、関数そのものを注入する）
	sleep func(time.Duration)
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

func (c *ConfigurableSleeper) Sleep() {

	// セットしたタイマーをsleepで定義した、SpyTimeのインスタンスメソッドのSleepを呼びだす
	c.sleep(c.duration)
}

const write = "write"
const sleep = "sleep"

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func (d DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

func (w MyWriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func TestCountdown(t *testing.T) {
	// *bytes.Bufferは内部でWriteメソッドを持っている
	// io.WriterのUSBの規格（現状のルール）＝インターフェース　を満たしている

	buffer := &bytes.Buffer{}
	spySleeper := &SpySleeper{}

	// ここにモックを差し込み、Sleepメソッドを自身のメソッドに置き換え、カウントする
	Countdown(buffer, spySleeper)

	got := buffer.String()
	want := `3
2
1
Go!`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

	if spySleeper.Calls != 4 {
		t.Errorf("not enough calls to sleeper, want 4 got %d", spySleeper.Calls)
	}

	t.Run("sleep before every print", func(t *testing.T) {
		spySleepPrinter := &CountdownOperationsSpy{}
		CountdownPrinter(spySleepPrinter, spySleepPrinter)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
		}
	})

}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second

	spyTime := &SpyTime{}
	//　タイマーを５秒セット、実行可能なコールバック関数を初期化
	sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
	// 実行可能な関数を呼び出す
	sleeper.Sleep()

	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
	}
}

// func Countdown(out *bytes.Buffer) {
// 	out.Write([]byte("3"))
// }

const finalWord = "Go!"
const countdownStart = 3

// ルールを制約
type Sleeper interface {
	Sleep()
}

// mock
type SpySleeper struct {
	// ここでモックを使い、何回呼ばれたカウントする
	Calls int
}

// 　Sleeperインターフェースの実装を満たす
// ここで自身のメソッドを差し込んで使用することができる
func (s *SpySleeper) Sleep() {
	s.Calls++
}

func Countdown(out io.Writer, time *SpySleeper) {
	for i := countdownStart; i > 0; i-- {
		time.Sleep()
		fmt.Fprintln(out, i)

	}
	time.Sleep()
	fmt.Fprint(out, finalWord)
}

func CountdownPrinter(out io.Writer, sleeper Sleeper) {

	// sleep
	for i := countdownStart; i > 0; i-- {
		// リストにsleepを追加
		sleeper.Sleep()
	}

	// 書き込みのテスト
	for i := countdownStart; i > 0; i-- {
		//　io.Writerで標準出力（bufに書き込み）
		fmt.Fprintln(out, i)

	}
	sleeper.Sleep()
	fmt.Fprint(out, finalWord)
}
