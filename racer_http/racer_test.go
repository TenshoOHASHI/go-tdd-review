package racerHttp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	slowURL := "http://www.facebook.com"
	fastURL := "http://www.quii.co.uk"

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// 早い方を返す
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration > bDuration {
		return a
	}
	return b

}

var tenSecondTimeout = 10 * time.Second

func TestMockRacer(t *testing.T) {
	// 登録した無名関数をリクエストのタイミングで発火し、内部でコールバック関数として実行する
	slowServer := makeDelayedServer(11 * time.Millisecond)
	fastServer := makeDelayedServer(9 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL

	got, err := RacerSelect(slowURL, fastURL, tenSecondTimeout)

	if err != nil {
		t.Errorf("got time out: %s", err)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

// 時間をラップしてあげる
// Dry-ing up（同じコードを重複せずに使用する）
func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

// 複数処理、またタイムアウトなどデフォルトで設定できる
func RacerSelect(a, b string, timeout time.Duration) (winner string, err error) {
	// チャンネルが閉じた方を先に読み出す
	// チャンネルが閉じるまで、待機状態になる
	select {
	// blocking
	case <-ping(a):
		return a, nil
	// blocking
	case <-ping(b):
		return b, nil
	// １０秒後にタイムアウトを設定
	case <-time.After(timeout):
		return "", fmt.Errorf("time out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan struct{} {
	// リクエストがあるたびに、空のチャンネルを作成
	// ここでmakeを使い、ゼロ値初期化させる, nilの場合、<-読み込み時に送信できずにパニックになってしまう
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		// チャンネルを閉じると、待機中のチャンネルから値が取り出される
		close(ch)
	}()

	return ch
}
