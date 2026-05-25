package main

import (
	"fmt"
	"io"
	"os"
)

func Countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}

func main() {
	// 値 -> os.Stdout -> 正体*os.File -> Writeメソッドを保持　＝　インターフェースのルールを満たしている
	Countdown(os.Stdout)
}
