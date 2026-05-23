package iteration

import (
	"strings"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := RepeatMyBuilder("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		// 引用符つきで、文字列リテラルとして安全に表示する
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func Repeat(character string, n int) string {
	var repeated string

	for range n {
		repeated = repeated + character
	}
	return repeated
}

func RepeatBuilder(character string, n int) string {
	var repeated strings.Builder
	repeated.Grow(len(character) * n)
	for range n {
		repeated.WriteString(character)
	}
	return repeated.String()
}

func RepeatMyBuilder(character string, n int) string {
	buf := make([]byte, 0, len(character)*n)

	for range n {
		// character[0]をインデックスでアクセスすると、バイトを返す（文字列は読み取り専用）
		// appendの場合、展開する際は、内部でbyteに変換してくれている。byte(character)...
		buf = append(buf, character...)
	}
	// バイト列を文字列に変換
	return string(buf)
}

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 5)
	}
}

func BenchmarkRepeatBuilder(b *testing.B) {
	for b.Loop() {
		RepeatBuilder("a", 5)
	}
}

// 予め、容量を確保すると、速度が変わってくる
// builderの場合、Growで容量を事前に確保してあげるといい
// -benchmemで容量を確認できる
// 小さく単純な処理では、[]byte + 事前容量確保のほうが速いことがある。
// strings.Builder は文字列構築用の標準APIで、意図が明確で扱いやすい。
func BenchmarkRepeatMyBuilder(b *testing.B) {
	for b.Loop() {
		RepeatMyBuilder("a", 5)
	}
}
