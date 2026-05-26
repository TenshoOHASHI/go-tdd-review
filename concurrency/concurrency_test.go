package concurrency_test

import (
	"reflect"
	"testing"
	"time"
)

// メソッドが必要の場合は、 関数インジェクションを使用する
type WebsiteChecker func(string) bool

type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	// 確認ように、マップを作成する
	results := make(map[string]bool)
	resultChannel := make(chan result)

	// ループするたびに、新しいゴルチーンが立ち上がり、並列処理をする
	for _, url := range urls {

		// 匿名関数は宣言と同時に実行できる
		// グローバルスコープを参照し、内部で使用可能
		go func(u string) {
			//　urlを参照し、渡す必要がある（独立したコピーを持たせる）
			// DIで設定した関数を呼び出す
			// results[url] = wc(u)
			resultChannel <- result{u, wc(u)} //　結果をチャネルに送信する, 個々に送信
		}(url)
	}

	// マップにデータを格納する際に、個々にデータを取り出し、格納したあげる
	for range len(urls) {
		// 結果を取り出して、マップに格納する
		result := <-resultChannel
		// ここで複数の結果をマップに格納する
		results[result.string] = result.bool
	}
	return results
}

// WebsiteCheckerの関数インジェクションとして使用する
func mockWebsiteChecker(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}
}

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

// ベンチマークを使用
func BenchmarkChecker(b *testing.B) {
	urls := make([]string, 100)
	for i, _ := range urls {
		urls[i] = "a url" // インデックス番号を指定して、そこに文字列の固定値を入れる
	}

	// b.NをLoopメソッドに格納
	// 並列処理で、役21.28ミリ秒に短縮、約同期処理より１００倍はやい
	for b.Loop() {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}

	// time.Sleep(2 * time.Second)

}
