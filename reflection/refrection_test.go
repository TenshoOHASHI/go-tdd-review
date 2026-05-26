package reflection_test

import (
	"reflect"
	"testing"
)

// インターフェースとは、間違った型をコンパイラー時に弾くために使用する
// ただ、コンパイル時に型がわからない関数を書きたい時に、any型を使用するinterface{}
// any型の場合は、すべての値を受け入れるけど、安全性が担保されなくなる。コンパイル時にエラーを検知できない
func TestReflection(t *testing.T) {
	expected := "Chris"
	var got []string // spyで記録

	// 匿名構造体を用意
	x := struct {
		Name string
	}{expected}

	walk(x, func(input string) {
		got = append(got, input)
	})

	if len(got) != 1 {
		t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
	}

	if got[0] != expected {
		t.Errorf("got %q, want %q", got[0], expected)
	}
}

func walk(x any, f func(input string)) {
	f("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}

func walkReflection(x any, f func(input string)) {
	val := reflect.ValueOf(x)
	// field := val.Field(0) // フィールドの1番目の要素を取得＝Name

	for i := range val.NumField() {
		field := val.Field(i)

		// フィールドが文字列のみ追加
		if field.Kind() == reflect.String {
			f(field.String())
		}

		if field.Kind() == reflect.Struct {
			// 見つかった構造体を元の構造体の値に戻す
			// value 構造体発見 -> interface() -> Profile
			walkReflection(field.Interface(), f)
		}

	}
}

func walkRefactor(x any, fn func(input string)) {
	val := getValue(x)

	// フィールド数にアクセスする前に、スライスを
	if val.Kind() == reflect.Slice {
		for i := range val.Len() {
			// スライスのi番目の構造体を渡す(単独の構造体として取り出す)
			walkRefactor(val.Index(i).Interface(), fn)
		}
		return
	}

	// スライスが渡ってくると、フィールド数がなくパニックになる
	for i := range val.NumField() {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walkRefactor(field.Interface(), fn)
		}
	}
}

func walkRefactorSwitch(x any, fn func(input string)) {
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		walkRefactorSwitch(value.Interface(), fn)
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())

	case reflect.Struct:
		for i := range val.NumField() {
			walkValue(val.Field(i))
		}

	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			walkValue(val.Index(i))

		}

	// Goの場合は順序を保証しないため、順番が間違えると、エラーが生じる
	case reflect.Map:
		// キーが存在しない場合エラー発生する、そのためキーを取り出し、存在するキーのみ、フィールドから個別に取り出す
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}

	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walkRefactorSwitch(v.Interface(), fn)
		}

	case reflect.Func:
		// 引数はなしで、関数を呼び出す、複数の戻り値をリストに格納
		valFnResult := val.Call(nil)
		// リストから個々のreflect.valueを取り出す
		for _, res := range valFnResult {
			walkRefactorSwitch(res.Interface(), fn)
		}
	}
}

func getValue(x any) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		// ポインター型の場合は、通常の構造体に戻す必要がある
		val = val.Elem()
	}
	return val
}

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{} // 入力
		ExpectedCalls []string    //　結果
	}{
		// 一番外側はデータの集合リスト
		// 内側は個々の構造体を初期化
		{"Struct with one string filed",
			//　構造体で渡す理由は、箱を用意して、reflectで箱を読ませるため
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},

		{"Struct with none string filed",
			//　構造体で渡す理由は、箱を用意して、reflectで箱を読ませるため
			struct {
				Name string
				Age  int
			}{"Chris", 0},
			[]string{"Chris"},
		},

		{"Struct with multi string filed",
			//　構造体で渡す理由は、箱を用意して、reflectで箱を読ませるため
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},

		{"Nested fields",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},

		{"Slices",
			[]Profile{
				{33, "London"},
				{32, "Tokyo"},
			},
			[]string{"London", "Tokyo"},
		},

		{
			"Arrays",
			[2]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},

		// {
		// 	// 順序を保証しない
		// 	"Maps",
		// 	map[string]string{
		// 		"Foo": "Bar",
		// 		"Baz": "Boz",
		// 	},
		// 	[]string{"Bar", "Boz"},
		// },
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walkRefactorSwitch(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("with map", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}
		var got []string
		walkRefactorSwitch(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})

	t.Run("with channel", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Berlin"}
			aChannel <- Profile{34, "Okinawa"}
			close(aChannel)
		}()

		var got []string
		want := []string{"Berlin", "Okinawa"}

		walkRefactorSwitch(aChannel, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("with function", func(t *testing.T) {
		aFunction := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Katowice"}
		}

		var got []string
		want := []string{"Berlin", "Katowice"}

		walkRefactorSwitch(aFunction, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

func assertContains(t *testing.T, haystack []string, needle string) {
	t.Helper()
	contains := false

	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
	}
}
