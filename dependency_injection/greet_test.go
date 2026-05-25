package dependencyinjection_test

import (
	"bytes"
	"fmt"
	"testing"
)

// type Greet struct {
// 	buf *bytes.Buffer
// 	str string
// }

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Greet(writer *bytes.Buffer, name string) {
	// 文字列をbuffに書き込んでいる
	fmt.Fprintf(writer, "Hello, %s", name)
}
