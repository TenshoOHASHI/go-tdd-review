package context_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// リクエストからコンテキスト内容を受け取る
		ctx := r.Context()

		// データの送受信用にチャンネルを作成
		data := make(chan string, 1)
		// 非同期処理でデータを取得
		go func() {
			data <- store.Fetch()
		}()

		// どっちが先に呼ばれた方を優先するため、競合を回避できる
		select {
		// データを受け取り状態になるまで待機
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
		// 標準出力
		fmt.Fprint(w, store.Fetch())
	}
}

func Server2(store Store2) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 今回はリクエストを受け取って、キャンセルされた後の、挙動、データの書き込みを確認するため、非同期処理で対応
		data, err := store.Fetch(r.Context())
		if err != nil {
			return
		}

		fmt.Fprint(w, data)
	}
}

type Store interface {
	// fetchする際に、親のコンテキストを作成
	Fetch() string
	Cancel()
}

type Store2 interface {
	// fetchする際に、親のコンテキストを作成
	Fetch(ctx context.Context) (string, error)
	Cancel()
}

type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	// 非同期でデータを受信
	go func() {
		// result := make([]rune, 0, len([]rune(s.response)))
		var result strings.Builder

		// 複数のデータ受信を１文字ずつ処理
		for _, c := range s.response {
			select {
			// もしキャンセルがあった場合、受信用のチャンネルが閉じ、終了する
			// Doneで受信用のチャンネルを取得する、そして待機状態になる
			case <-ctx.Done():

				s.t.Log("spy store got cancel")
				return
			//　そうでない場合は、
			case <-time.After(10 * time.Millisecond):
				result.WriteRune(c)
			}
		}
		// データを受け取るたびに、データに送信
		data <- result.String()
	}()

	// 外側の処理も同時に処理を開始：キャンセルが先かデータが先か
	select {
	// もしキャンセルされていたら、ここでエラーを出力
	case <-ctx.Done():
		s.Cancel()
		return "", ctx.Err()
	// データがあれば、それを取り出す
	case res := <-data:
		return res, nil
	}

}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

type StubStore struct {
	response string
}

// "hello, world"がかえってくる
func (s *StubStore) Fetch() string {
	return s.response
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("store was told to cancel")
	}
}

func TestServer(t *testing.T) {
	data := "hello, world"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server2(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}

		store.assertWasNotCancelled()
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server2(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		// キャンセル可能な 親のcontext を作る
		//　受信用のDone()チャンネルを取得する　 <- ctx.Done()でまつ
		cancellingCtx, cancel := context.WithCancel(request.Context())
		// 5ミリ秒後に cancel() を呼ぶと子のコンテキストがキャンセルされる（Done() channel, cancel()）
		//  内部の done channel が close される
		// <-ctx.Done() の待機が解除される
		time.AfterFunc(5*time.Millisecond, cancel)
		// request にキャンセル可能な context をセットする。
		request = request.WithContext(cancellingCtx)

		// response := httptest.NewRecorder()
		response := &SpyResponseWriter{}

		svr.ServeHTTP(response, request)

		// キャンセルした時に、中身に何も記載されているないことを確認する
		if response.written {
			t.Error("a response should not have been written")
		}

		store.assertWasCancelled()
	})
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}
