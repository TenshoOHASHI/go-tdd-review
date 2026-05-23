package map_test

import (
	"errors"
	"testing"
)

// Dictionary は単語をキー、説明文を値として持つ独自の map 型。
// map[string]string を元にしているので、Search や Add のようなメソッドを追加できる。
type Dictionary map[string]string

var (
	ErrNotFound         = DictionaryErr("could not find the word you where looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

// エラーを再利用可能にし、不変にする
type DictionaryErr string

// 初期化したエラーメッセージを文字列で返す
func (e DictionaryErr) Error() string {
	return string(e)
}

// Search は Dictionary から word に対応する説明文を探す。
// 見つからない場合は ErrNotFound を返す。
func (d Dictionary) Search(word string) (string, error) {
	// map の2値受け取りで、値と存在有無を同時に取得する。
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	return definition, nil
}

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, err := dictionary.Search("test")
		want := "this is just a test"
		assertNoError(t, err)
		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		assertError(t, err, ErrNotFound)
	})

}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	// エラーを期待しているので、nil ならここでテストを止める。
	if got == nil {
		t.Fatal("expected to get an error")
	}

	// errors.Is は、将来エラーが fmt.Errorf("%w", err) で包まれても判定できる。
	if !errors.Is(got, want) {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected no error, got %q", got)
	}
}

// Add は未登録の word だけ Dictionary に追加する。
// map は内部データへの参照を持つため、ポインタレシーバーでなくても中身を変更できる。
func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	switch {
	// まだ単語が存在しないので、追加する
	case errors.Is(err, ErrNotFound):
		d[word] = definition
		return nil
	//　値が存在するので、エラーを返す
	case err == nil:
		return ErrWordExists
	default:
		return err
	}
}

// Write the test first
func TestAdd(t *testing.T) {

	t.Run("new word", func(t *testing.T) {
		//　varで宣言すると、内部のハッシュテーブルがまだ作成されず、nilになってしますため、避ける
		// var dictionary Dictionary はNG
		dictionary := Dictionary{}
		word := "test"
		definition := "this is just a test"

		err := dictionary.Add(word, definition)
		assertNoError(t, err)
		assertDefinition(t, dictionary, word, definition)

	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		// すでにtestが入っているので、重複する
		err := dictionary.Add(word, "new test")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})

}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	// 単語が見つからない場合は、エラーを返す
	case ErrNotFound:
		return ErrWordDoesNotExist
	// 単語が存在する場合は、値を変更
	case nil:
		d[word] = definition
		return nil
	// それ以外のエラー
	default:
		return err
	}
}

// 存在しないキーを削除しても何も起きないから、シンプルに書く
func (d Dictionary) Delete(word string) {
	delete(d, word)
}

func assertDefinition(t *testing.T, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if definition != got {
		t.Errorf("got %q want %q", got, definition)
	}
}

// 既存のキーに対して値を変更する
func TestUpdate(t *testing.T) {

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		newDefinition := "new definition"
		dictionary := Dictionary{word: definition}

		err := dictionary.Update(word, newDefinition)

		assertNoError(t, err)
		assertDefinition(t, dictionary, word, newDefinition)
	})

	t.Run("empty word and definition", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}

		err := dictionary.Update(word, definition)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}

		dictionary.Delete(word)
		_, err := dictionary.Search(word)
		assertError(t, err, ErrNotFound)

	})
}
