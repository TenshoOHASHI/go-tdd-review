package sync_test

import (
	"sync"
	"testing"
)

func TestCount(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := NewCounter()
		wantedCount := 1000

		// down := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		// chanを使用する場合は、同期的に便利だけど、基本AからBに伝える時に、使用
		for range wantedCount {
			go func(w *sync.WaitGroup) {
				counter.Inc()
				defer w.Done()
				// defer close(down)
			}(&wg)
		}

		// <-down
		// メインゴルチーンが終わらないように待機する
		wg.Wait()
		// Mutexはコピーしてはいけない
		assertCounter(t, counter, 1000)

	})

}

func assertCounter(t *testing.T, got *Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}

// go　vetで微妙なバグを警告することができます。
func NewCounter() *Counter {
	return &Counter{}
}

type Counter struct {
	mu    sync.Mutex // *ここで直接埋め込むのは危険、公開時に呼び出し先から、mutex操作ができてします
	value int
}

func (c *Counter) Value() int {
	return c.value
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}
