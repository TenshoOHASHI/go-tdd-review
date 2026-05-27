package context

import (
	"context"
	"errors"
	"sync"
)

var ErrCanceled = errors.New("context canceled")

// 関数値の型として定義
type CancelFunc func()

type Context interface {
	Done() <-chan struct{}
	Err() error
}

type cancelCtx struct {
	done chan struct{}
	once sync.Once
	mu   sync.Mutex
	err  error
}

func (c *cancelCtx) Done() <-chan struct{} {
	if c.done == nil {
		c.done = make(chan struct{})
	}
	return c.done
}

func (c *cancelCtx) cancel() {
	c.mu.Lock()
	// キャンセル時に、グローバルエラー変数をインスタンス変数に格納する
	c.err = context.Canceled
	c.mu.Unlock()
	// done チャネルで受信待ちしている goroutine がブロック解除される
	close(c.done)
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func WithCancel(parent Context) (Context, CancelFunc) {
	// 1. 新しいコンテキストを作る
	ctx := &cancelCtx{
		done: make(chan struct{}), // 専用のチャネルを作る
	}

	// 2. 「キャンセルするための関数」を作る
	cancel := func() {
		ctx.cancel()
	}

	return ctx, cancel
}

// 使う側はこうです。

//   ctx, cancel := WithCancel(nil)

//   go func() {
//   	<-ctx.Done()
//   	fmt.Println(ctx.Err())
//   }()

//   cancel()
